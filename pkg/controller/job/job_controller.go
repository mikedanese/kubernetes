/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package job

import (
	"reflect"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client/cache"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/controller"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/controller/framework"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/workqueue"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
	"github.com/golang/glog"
)

type JobManager struct {
	kubeClient client.Interface
	podControl controller.PodControlInterface

	// An rc is temporarily suspended after creating/deleting these many jobs.
	// It resumes normal action after observing the watch events for them.
	burstJobs int
	// To allow injection for testing
	syncHandler func(jobKey string) error

	// A TTLCache of pod creates/deletes each rc expects to see
	expectations controller.ControllerExpectationsInterface
	// podStoreSynced returns true if the pod store has been synced at least once.
	// Added as a member to the struct to allow injection for testing.
	podStoreSynced func() bool

	// A store of jobs, populated by the jobController
	jobStore cache.StoreToJobLister
	// A store of pods, populated by the podController
	podStore cache.StoreToPodLister
	// Watches changes to all jobs
	jobController *framework.Controller
	// Watches changes to all pods
	podController *framework.Controller
	// Controllers that need to be updated
	queue *workqueue.Type
}

var (
	burstJobs        = 500
	fullResyncPeriod = 30 * time.Second
	podRelistPeriod  = 5 * time.Minute
)

func NewJobManager(kubeClient client.Interface, burstJobs int) *JobManager {
	jm := &JobManager{
		kubeClient:   kubeClient,
		burstJobs:    burstJobs,
		expectations: controller.NewControllerExpectations(),
		queue:        workqueue.New(),
	}

	jm.jobStore.Store, jm.jobController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func() (runtime.Object, error) {
				return jm.kubeClient.Jobs(api.NamespaceAll).List(labels.Everything())
			},
			WatchFunc: func(rv string) (watch.Interface, error) {
				return jm.kubeClient.Jobs(api.NamespaceAll).Watch(labels.Everything(), fields.Everything(), rv)
			},
		},
		&api.Job{},
		fullResyncPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc: jm.enqueueJob,
			UpdateFunc: func(old, cur interface{}) {
				jm.enqueueJob(cur)
			},
			// This will enter the sync loop and no-op, becuase the controller has been deleted from the store.
			// Note that deleting a controller immediately after scaling it to 0 will not work. The recommended
			// way of achieving this is by performing a `stop` operation on the controller.
			DeleteFunc: jm.enqueueJob,
		},
	)

	jm.podStore.Store, jm.podController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc: func() (runtime.Object, error) {
				return jm.kubeClient.Pods(api.NamespaceAll).List(labels.Everything(), fields.Everything())
			},
			WatchFunc: func(rv string) (watch.Interface, error) {
				return jm.kubeClient.Pods(api.NamespaceAll).Watch(labels.Everything(), fields.Everything(), rv)
			},
		},
		&api.Pod{},
		podRelistPeriod,
		framework.ResourceEventHandlerFuncs{
			AddFunc:    jm.addPod,
			UpdateFunc: jm.updatePod,
			DeleteFunc: jm.deletePod,
		},
	)

	jm.syncHandler = jm.syncJob
	jm.podStoreSynced = jm.podController.HasSynced

	return jm
}

func (jm *JobManager) Run(workers int, stopCh <-chan struct{}) {
	defer util.HandleCrash()
	go jm.jobController.Run(stopCh)
	go jm.podController.Run(stopCh)
	for i := 0; i < workers; i++ {
		go util.Until(jm.worker, time.Second, stopCh)
	}
	<-stopCh
}

func (jm *JobManager) worker() {
	for {
		func() {
			key, quit := jm.queue.Get()
			if quit {
				return
			}
			defer jm.queue.Done(key)
			err := jm.syncHandler(key.(string))
			if err != nil {
				glog.Errorf("Error syncing replication controller: %v", err)
			}
		}()
	}
}

func (jm *JobManager) syncJob(key string) error {
	glog.Infof("syncJob: %#v", key)
	defer func(startTime time.Time) {
		glog.V(4).Infof("Finished syncing job %q (%v)", key, time.Now().Sub(startTime))
	}(time.Now())

	obj, exists, err := jm.jobStore.Store.GetByKey(key)
	if !exists {
		glog.Infof("job %v does not exist")
		return nil
	}
	if err != nil {
		glog.Infof("err retrieving job %q from store: %v", key, err)
		jm.queue.Add(key)
		return err
	}
	job := *obj.(*api.Job)
	if !jm.podStoreSynced() {
		time.Sleep(30 * time.Second)
		jm.enqueueJob(job)
	}

	// Check the expectations of the rc before counting active pods, otherwise a new pod can sneak in
	// and update the expectations after we've retrieved active pods from the store. If a new pod enters
	// the store after we've checked the expectation, the rc sync is just deferred till the next relist.
	if err != nil {
		glog.Errorf("Couldn't get key for job %#v: %v", job, err)
		return err
	}
	jobNeedsSync := jm.expectations.SatisfiedExpectations(key)
	podList, err := jm.podStore.Pods(job.Namespace).List(labels.Set(job.Spec.Selector).AsSelector())
	if err != nil {
		glog.Errorf("Error getting pods for rc %q: %v", key, err)
		jm.queue.Add(key)
		return err
	}

	// TODO: Do this in a single pass, or use an index.
	if jobNeedsSync {
		jm.updatePodJobs(podList.Items, &job)
	}

	// Always updates status as pods come up or die.
	if err := updateStatus(jm.kubeClient.Jobs(job.Namespace), job, podList.Items); err != nil {
		glog.V(2).Infof("Failed to update replica count for job %v, requeuing", job.Name)
		jm.enqueueJob(&job)
	}

	return nil
}

func (jm *JobManager) updatePodJobs(pods []api.Pod, job *api.Job) {
	glog.Infof("updatePodJobs: %#v", jm.expectations)
	diff := len(filteredPods) - rc.Spec.Replicas
	rcKey, err := controller.KeyFunc(rc)
	if err != nil {
		glog.Errorf("Couldn't get key for replication controller %#v: %v", rc, err)
		return
	}
	if diff < 0 {
		diff *= -1
		if diff > rm.burstReplicas {
			diff = rm.burstReplicas
		}
		rm.expectations.ExpectCreations(rcKey, diff)
		wait := sync.WaitGroup{}
		wait.Add(diff)
		glog.V(2).Infof("Too few %q/%q replicas, need %d, creating %d", rc.Namespace, rc.Name, rc.Spec.Replicas, diff)
		for i := 0; i < diff; i++ {
			go func() {
				defer wait.Done()
				if err := rm.podControl.CreateReplica(rc.Namespace, rc); err != nil {
					// Decrement the expected number of creates because the informer won't observe this pod
					glog.V(2).Infof("Failed creation, decrementing expectations for controller %q/%q", rc.Namespace, rc.Name)
					rm.expectations.CreationObserved(rcKey)
					util.HandleError(err)
				}
			}()
		}
		wait.Wait()
	} else if diff > 0 {
		if diff > rm.burstReplicas {
			diff = rm.burstReplicas
		}
		rm.expectations.ExpectDeletions(rcKey, diff)
		glog.V(2).Infof("Too many %q/%q replicas, need %d, deleting %d", rc.Namespace, rc.Name, rc.Spec.Replicas, diff)
		// No need to sort pods if we are about to delete all of them
		if rc.Spec.Replicas != 0 {
			// Sort the pods in the order such that not-ready < ready, unscheduled
			// < scheduled, and pending < running. This ensures that we delete pods
			// in the earlier stages whenever possible.
			sort.Sort(controller.ActivePods(filteredPods))
		}

		wait := sync.WaitGroup{}
		wait.Add(diff)
		for i := 0; i < diff; i++ {
			go func(ix int) {
				defer wait.Done()
				if err := rm.podControl.DeletePod(rc.Namespace, filteredPods[ix].Name); err != nil {
					// Decrement the expected number of deletes because the informer won't observe this deletion
					glog.V(2).Infof("Failed deletion, decrementing expectations for controller %q/%q", rc.Namespace, rc.Name)
					rm.expectations.DeletionObserved(rcKey)
				}
			}(i)
		}
		wait.Wait()
	}

}

func (jm *JobManager) getPodJob(pod *api.Pod) *api.Job {
	glog.Infof("getPodJob: %#v", pod)
	jobs, err := jm.jobStore.GetPodJobs(pod)
	if err != nil {
		glog.V(4).Infof("unable to find job for pod %v: %v", pod, err)
		return nil
	}
	return getFirstCreated(jobs)
}

func (jm *JobManager) enqueueJob(obj interface{}) {
	glog.Infof("enqueJob: %#v", obj)
	key, err := controller.KeyFunc(obj)
	if err != nil {
		glog.Errorf("Couldn't get key for object %+v: %v", obj, err)
		return
	}
	jm.queue.Add(key)
}

func (jm *JobManager) addPod(obj interface{}) {
	glog.Infof("addJob: %#v", obj)
	pod := obj.(*api.Pod)
	if job := jm.getPodJob(pod); job != nil {
		jobKey, err := controller.KeyFunc(job)
		if err != nil {
			glog.Errorf("Couldn't get key for job %#v: %v", job, err)
			return
		}
		jm.expectations.CreationObserved(jobKey)
		jm.enqueueJob(job)
	}
}

func (jm *JobManager) updatePod(oldObj, newObj interface{}) {
	glog.Infof("updateJob: %#v", newObj)
	if api.Semantic.DeepEqual(oldObj, newObj) {
		// A periodic relist will send update events for all known pods.
		return
	}

	newPod := newObj.(*api.Pod)
	if job := jm.getPodJob(newPod); job != nil {
		jm.enqueueJob(job)
	}
	oldPod := oldObj.(*api.Pod)

	if !reflect.DeepEqual(newPod.Labels, oldPod.Labels) {
		if oldJob := jm.getPodJob(oldPod); oldJob != nil {
			jm.enqueueJob(oldJob)
		}
	}
}

func (jm *JobManager) deletePod(obj interface{}) {
	glog.Infof("deleteJob: %#v", obj)
	pod, ok := obj.(*api.Pod)

	// When a delete is dropped, the relist will notice a pod in the store not
	// in the list, leading to the insertion of a tombstone object which contains
	// the deleted key/value. Note that this value might be stale. If the pod
	// changed labels the new rc will not be woken up till the periodic resync.
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			glog.Errorf("Couldn't get object from tombstone %+v, could take up to %v before a controller recreates a replica", obj, controller.ExpectationsTimeout)
			return
		}
		pod, ok = tombstone.Obj.(*api.Pod)
		if !ok {
			glog.Errorf("Tombstone contained object that is not a pod %+v, could take up to %v before controller recreates a replica", obj, controller.ExpectationsTimeout)
			return
		}
	}
	if job := jm.getPodJob(pod); job != nil {
		jobKey, err := controller.KeyFunc(job)
		if err != nil {
			glog.Errorf("Couldn't get key for job %#v: %v", job, err)
			return
		}
		jm.expectations.DeletionObserved(jobKey)
		jm.enqueueJob(job)
	}
}

func updateStatus(jobs client.JobInterface, job api.Job, pods []api.Pod) error {
	return nil
}

func getFirstCreated(jobs []*api.Job) *api.Job {
	return jobs[0]
}
