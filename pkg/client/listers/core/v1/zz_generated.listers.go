/*
Copyright 2016 The Kubernetes Authors.

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

// This file was automatically generated by lister-gen with arguments: --input-dirs=[k8s.io/kubernetes/pkg/api,k8s.io/kubernetes/pkg/api/v1,k8s.io/kubernetes/pkg/apis/abac,k8s.io/kubernetes/pkg/apis/abac/v0,k8s.io/kubernetes/pkg/apis/abac/v1beta1,k8s.io/kubernetes/pkg/apis/apps,k8s.io/kubernetes/pkg/apis/apps/v1alpha1,k8s.io/kubernetes/pkg/apis/authentication,k8s.io/kubernetes/pkg/apis/authentication/v1beta1,k8s.io/kubernetes/pkg/apis/authorization,k8s.io/kubernetes/pkg/apis/authorization/v1beta1,k8s.io/kubernetes/pkg/apis/autoscaling,k8s.io/kubernetes/pkg/apis/autoscaling/v1,k8s.io/kubernetes/pkg/apis/batch,k8s.io/kubernetes/pkg/apis/batch/v1,k8s.io/kubernetes/pkg/apis/batch/v2alpha1,k8s.io/kubernetes/pkg/apis/certificates,k8s.io/kubernetes/pkg/apis/certificates/v1alpha1,k8s.io/kubernetes/pkg/apis/componentconfig,k8s.io/kubernetes/pkg/apis/componentconfig/v1alpha1,k8s.io/kubernetes/pkg/apis/extensions,k8s.io/kubernetes/pkg/apis/extensions/v1beta1,k8s.io/kubernetes/pkg/apis/imagepolicy,k8s.io/kubernetes/pkg/apis/imagepolicy/v1alpha1,k8s.io/kubernetes/pkg/apis/policy,k8s.io/kubernetes/pkg/apis/policy/v1alpha1,k8s.io/kubernetes/pkg/apis/rbac,k8s.io/kubernetes/pkg/apis/rbac/v1alpha1,k8s.io/kubernetes/pkg/apis/storage,k8s.io/kubernetes/pkg/apis/storage/v1beta1]

package v1

import (
	api "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	v1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
)

// ComponentStatusLister helps list ComponentStatuses.
type ComponentStatusLister interface {
	// List lists all ComponentStatuses in the indexer.
	List(selector labels.Selector) (ret []*v1.ComponentStatus, err error)
	// Get retrieves the ComponentStatus from the index for a given name.
	Get(name string) (*v1.ComponentStatus, error)
}

// componentStatusLister implements the ComponentStatusLister interface.
type componentStatusLister struct {
	indexer cache.Indexer
}

// NewComponentStatusLister returns a new ComponentStatusLister.
func NewComponentStatusLister(indexer cache.Indexer) ComponentStatusLister {
	return &componentStatusLister{indexer: indexer}
}

// List lists all ComponentStatuses in the indexer.
func (s *componentStatusLister) List(selector labels.Selector) (ret []*v1.ComponentStatus, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ComponentStatus))
	})
	return ret, err
}

// Get retrieves the ComponentStatus from the index for a given name.
func (s *componentStatusLister) Get(name string) (*v1.ComponentStatus, error) {
	key := &v1.ComponentStatus{ObjectMeta: v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("componentstatus"), name)
	}
	return obj.(*v1.ComponentStatus), nil
}

// ConfigMapLister helps list ConfigMaps.
type ConfigMapLister interface {
	// List lists all ConfigMaps in the indexer.
	List(selector labels.Selector) (ret []*v1.ConfigMap, err error)
	// ConfigMaps returns an object that can list and get ConfigMaps.
	ConfigMaps(namespace string) ConfigMapNamespaceLister
}

// configMapLister implements the ConfigMapLister interface.
type configMapLister struct {
	indexer cache.Indexer
}

// NewConfigMapLister returns a new ConfigMapLister.
func NewConfigMapLister(indexer cache.Indexer) ConfigMapLister {
	return &configMapLister{indexer: indexer}
}

// List lists all ConfigMaps in the indexer.
func (s *configMapLister) List(selector labels.Selector) (ret []*v1.ConfigMap, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ConfigMap))
	})
	return ret, err
}

// ConfigMaps returns an object that can list and get ConfigMaps.
func (s *configMapLister) ConfigMaps(namespace string) ConfigMapNamespaceLister {
	return configMapNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ConfigMapNamespaceLister helps list and get ConfigMaps.
type ConfigMapNamespaceLister interface {
	// List lists all ConfigMaps in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ConfigMap, err error)
	// Get retrieves the ConfigMap from the indexer for a given namespace and name.
	Get(name string) (*v1.ConfigMap, error)
}

// configMapNamespaceLister implements the ConfigMapNamespaceLister
// interface.
type configMapNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ConfigMaps in the indexer for a given namespace.
func (s configMapNamespaceLister) List(selector labels.Selector) (ret []*v1.ConfigMap, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ConfigMap))
	})
	return ret, err
}

// Get retrieves the ConfigMap from the indexer for a given namespace and name.
func (s configMapNamespaceLister) Get(name string) (*v1.ConfigMap, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("configmap"), name)
	}
	return obj.(*v1.ConfigMap), nil
}

// EndpointsLister helps list Endpoints.
type EndpointsLister interface {
	// List lists all Endpoints in the indexer.
	List(selector labels.Selector) (ret []*v1.Endpoints, err error)
	// Endpoints returns an object that can list and get Endpoints.
	Endpoints(namespace string) EndpointsNamespaceLister
}

// endpointsLister implements the EndpointsLister interface.
type endpointsLister struct {
	indexer cache.Indexer
}

// NewEndpointsLister returns a new EndpointsLister.
func NewEndpointsLister(indexer cache.Indexer) EndpointsLister {
	return &endpointsLister{indexer: indexer}
}

// List lists all Endpoints in the indexer.
func (s *endpointsLister) List(selector labels.Selector) (ret []*v1.Endpoints, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Endpoints))
	})
	return ret, err
}

// Endpoints returns an object that can list and get Endpoints.
func (s *endpointsLister) Endpoints(namespace string) EndpointsNamespaceLister {
	return endpointsNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EndpointsNamespaceLister helps list and get Endpoints.
type EndpointsNamespaceLister interface {
	// List lists all Endpoints in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Endpoints, err error)
	// Get retrieves the Endpoints from the indexer for a given namespace and name.
	Get(name string) (*v1.Endpoints, error)
}

// endpointsNamespaceLister implements the EndpointsNamespaceLister
// interface.
type endpointsNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Endpoints in the indexer for a given namespace.
func (s endpointsNamespaceLister) List(selector labels.Selector) (ret []*v1.Endpoints, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Endpoints))
	})
	return ret, err
}

// Get retrieves the Endpoints from the indexer for a given namespace and name.
func (s endpointsNamespaceLister) Get(name string) (*v1.Endpoints, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("endpoints"), name)
	}
	return obj.(*v1.Endpoints), nil
}

// EventLister helps list Events.
type EventLister interface {
	// List lists all Events in the indexer.
	List(selector labels.Selector) (ret []*v1.Event, err error)
	// Events returns an object that can list and get Events.
	Events(namespace string) EventNamespaceLister
}

// eventLister implements the EventLister interface.
type eventLister struct {
	indexer cache.Indexer
}

// NewEventLister returns a new EventLister.
func NewEventLister(indexer cache.Indexer) EventLister {
	return &eventLister{indexer: indexer}
}

// List lists all Events in the indexer.
func (s *eventLister) List(selector labels.Selector) (ret []*v1.Event, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Event))
	})
	return ret, err
}

// Events returns an object that can list and get Events.
func (s *eventLister) Events(namespace string) EventNamespaceLister {
	return eventNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EventNamespaceLister helps list and get Events.
type EventNamespaceLister interface {
	// List lists all Events in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Event, err error)
	// Get retrieves the Event from the indexer for a given namespace and name.
	Get(name string) (*v1.Event, error)
}

// eventNamespaceLister implements the EventNamespaceLister
// interface.
type eventNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Events in the indexer for a given namespace.
func (s eventNamespaceLister) List(selector labels.Selector) (ret []*v1.Event, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Event))
	})
	return ret, err
}

// Get retrieves the Event from the indexer for a given namespace and name.
func (s eventNamespaceLister) Get(name string) (*v1.Event, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("event"), name)
	}
	return obj.(*v1.Event), nil
}

// LimitRangeLister helps list LimitRanges.
type LimitRangeLister interface {
	// List lists all LimitRanges in the indexer.
	List(selector labels.Selector) (ret []*v1.LimitRange, err error)
	// LimitRanges returns an object that can list and get LimitRanges.
	LimitRanges(namespace string) LimitRangeNamespaceLister
}

// limitRangeLister implements the LimitRangeLister interface.
type limitRangeLister struct {
	indexer cache.Indexer
}

// NewLimitRangeLister returns a new LimitRangeLister.
func NewLimitRangeLister(indexer cache.Indexer) LimitRangeLister {
	return &limitRangeLister{indexer: indexer}
}

// List lists all LimitRanges in the indexer.
func (s *limitRangeLister) List(selector labels.Selector) (ret []*v1.LimitRange, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.LimitRange))
	})
	return ret, err
}

// LimitRanges returns an object that can list and get LimitRanges.
func (s *limitRangeLister) LimitRanges(namespace string) LimitRangeNamespaceLister {
	return limitRangeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// LimitRangeNamespaceLister helps list and get LimitRanges.
type LimitRangeNamespaceLister interface {
	// List lists all LimitRanges in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.LimitRange, err error)
	// Get retrieves the LimitRange from the indexer for a given namespace and name.
	Get(name string) (*v1.LimitRange, error)
}

// limitRangeNamespaceLister implements the LimitRangeNamespaceLister
// interface.
type limitRangeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all LimitRanges in the indexer for a given namespace.
func (s limitRangeNamespaceLister) List(selector labels.Selector) (ret []*v1.LimitRange, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.LimitRange))
	})
	return ret, err
}

// Get retrieves the LimitRange from the indexer for a given namespace and name.
func (s limitRangeNamespaceLister) Get(name string) (*v1.LimitRange, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("limitrange"), name)
	}
	return obj.(*v1.LimitRange), nil
}

// NamespaceLister helps list Namespaces.
type NamespaceLister interface {
	// List lists all Namespaces in the indexer.
	List(selector labels.Selector) (ret []*v1.Namespace, err error)
	// Get retrieves the Namespace from the index for a given name.
	Get(name string) (*v1.Namespace, error)
}

// namespaceLister implements the NamespaceLister interface.
type namespaceLister struct {
	indexer cache.Indexer
}

// NewNamespaceLister returns a new NamespaceLister.
func NewNamespaceLister(indexer cache.Indexer) NamespaceLister {
	return &namespaceLister{indexer: indexer}
}

// List lists all Namespaces in the indexer.
func (s *namespaceLister) List(selector labels.Selector) (ret []*v1.Namespace, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Namespace))
	})
	return ret, err
}

// Get retrieves the Namespace from the index for a given name.
func (s *namespaceLister) Get(name string) (*v1.Namespace, error) {
	key := &v1.Namespace{ObjectMeta: v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("namespace"), name)
	}
	return obj.(*v1.Namespace), nil
}

// NodeLister helps list Nodes.
type NodeLister interface {
	// List lists all Nodes in the indexer.
	List(selector labels.Selector) (ret []*v1.Node, err error)
	// Get retrieves the Node from the index for a given name.
	Get(name string) (*v1.Node, error)
}

// nodeLister implements the NodeLister interface.
type nodeLister struct {
	indexer cache.Indexer
}

// NewNodeLister returns a new NodeLister.
func NewNodeLister(indexer cache.Indexer) NodeLister {
	return &nodeLister{indexer: indexer}
}

// List lists all Nodes in the indexer.
func (s *nodeLister) List(selector labels.Selector) (ret []*v1.Node, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Node))
	})
	return ret, err
}

// Get retrieves the Node from the index for a given name.
func (s *nodeLister) Get(name string) (*v1.Node, error) {
	key := &v1.Node{ObjectMeta: v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("node"), name)
	}
	return obj.(*v1.Node), nil
}

// PersistentVolumeLister helps list PersistentVolumes.
type PersistentVolumeLister interface {
	// List lists all PersistentVolumes in the indexer.
	List(selector labels.Selector) (ret []*v1.PersistentVolume, err error)
	// Get retrieves the PersistentVolume from the index for a given name.
	Get(name string) (*v1.PersistentVolume, error)
}

// persistentVolumeLister implements the PersistentVolumeLister interface.
type persistentVolumeLister struct {
	indexer cache.Indexer
}

// NewPersistentVolumeLister returns a new PersistentVolumeLister.
func NewPersistentVolumeLister(indexer cache.Indexer) PersistentVolumeLister {
	return &persistentVolumeLister{indexer: indexer}
}

// List lists all PersistentVolumes in the indexer.
func (s *persistentVolumeLister) List(selector labels.Selector) (ret []*v1.PersistentVolume, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PersistentVolume))
	})
	return ret, err
}

// Get retrieves the PersistentVolume from the index for a given name.
func (s *persistentVolumeLister) Get(name string) (*v1.PersistentVolume, error) {
	key := &v1.PersistentVolume{ObjectMeta: v1.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("persistentvolume"), name)
	}
	return obj.(*v1.PersistentVolume), nil
}

// PersistentVolumeClaimLister helps list PersistentVolumeClaims.
type PersistentVolumeClaimLister interface {
	// List lists all PersistentVolumeClaims in the indexer.
	List(selector labels.Selector) (ret []*v1.PersistentVolumeClaim, err error)
	// PersistentVolumeClaims returns an object that can list and get PersistentVolumeClaims.
	PersistentVolumeClaims(namespace string) PersistentVolumeClaimNamespaceLister
}

// persistentVolumeClaimLister implements the PersistentVolumeClaimLister interface.
type persistentVolumeClaimLister struct {
	indexer cache.Indexer
}

// NewPersistentVolumeClaimLister returns a new PersistentVolumeClaimLister.
func NewPersistentVolumeClaimLister(indexer cache.Indexer) PersistentVolumeClaimLister {
	return &persistentVolumeClaimLister{indexer: indexer}
}

// List lists all PersistentVolumeClaims in the indexer.
func (s *persistentVolumeClaimLister) List(selector labels.Selector) (ret []*v1.PersistentVolumeClaim, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PersistentVolumeClaim))
	})
	return ret, err
}

// PersistentVolumeClaims returns an object that can list and get PersistentVolumeClaims.
func (s *persistentVolumeClaimLister) PersistentVolumeClaims(namespace string) PersistentVolumeClaimNamespaceLister {
	return persistentVolumeClaimNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PersistentVolumeClaimNamespaceLister helps list and get PersistentVolumeClaims.
type PersistentVolumeClaimNamespaceLister interface {
	// List lists all PersistentVolumeClaims in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.PersistentVolumeClaim, err error)
	// Get retrieves the PersistentVolumeClaim from the indexer for a given namespace and name.
	Get(name string) (*v1.PersistentVolumeClaim, error)
}

// persistentVolumeClaimNamespaceLister implements the PersistentVolumeClaimNamespaceLister
// interface.
type persistentVolumeClaimNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PersistentVolumeClaims in the indexer for a given namespace.
func (s persistentVolumeClaimNamespaceLister) List(selector labels.Selector) (ret []*v1.PersistentVolumeClaim, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PersistentVolumeClaim))
	})
	return ret, err
}

// Get retrieves the PersistentVolumeClaim from the indexer for a given namespace and name.
func (s persistentVolumeClaimNamespaceLister) Get(name string) (*v1.PersistentVolumeClaim, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("persistentvolumeclaim"), name)
	}
	return obj.(*v1.PersistentVolumeClaim), nil
}

// PodLister helps list Pods.
type PodLister interface {
	// List lists all Pods in the indexer.
	List(selector labels.Selector) (ret []*v1.Pod, err error)
	// Pods returns an object that can list and get Pods.
	Pods(namespace string) PodNamespaceLister
}

// podLister implements the PodLister interface.
type podLister struct {
	indexer cache.Indexer
}

// NewPodLister returns a new PodLister.
func NewPodLister(indexer cache.Indexer) PodLister {
	return &podLister{indexer: indexer}
}

// List lists all Pods in the indexer.
func (s *podLister) List(selector labels.Selector) (ret []*v1.Pod, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Pod))
	})
	return ret, err
}

// Pods returns an object that can list and get Pods.
func (s *podLister) Pods(namespace string) PodNamespaceLister {
	return podNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodNamespaceLister helps list and get Pods.
type PodNamespaceLister interface {
	// List lists all Pods in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Pod, err error)
	// Get retrieves the Pod from the indexer for a given namespace and name.
	Get(name string) (*v1.Pod, error)
}

// podNamespaceLister implements the PodNamespaceLister
// interface.
type podNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Pods in the indexer for a given namespace.
func (s podNamespaceLister) List(selector labels.Selector) (ret []*v1.Pod, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Pod))
	})
	return ret, err
}

// Get retrieves the Pod from the indexer for a given namespace and name.
func (s podNamespaceLister) Get(name string) (*v1.Pod, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("pod"), name)
	}
	return obj.(*v1.Pod), nil
}

// PodTemplateLister helps list PodTemplates.
type PodTemplateLister interface {
	// List lists all PodTemplates in the indexer.
	List(selector labels.Selector) (ret []*v1.PodTemplate, err error)
	// PodTemplates returns an object that can list and get PodTemplates.
	PodTemplates(namespace string) PodTemplateNamespaceLister
}

// podTemplateLister implements the PodTemplateLister interface.
type podTemplateLister struct {
	indexer cache.Indexer
}

// NewPodTemplateLister returns a new PodTemplateLister.
func NewPodTemplateLister(indexer cache.Indexer) PodTemplateLister {
	return &podTemplateLister{indexer: indexer}
}

// List lists all PodTemplates in the indexer.
func (s *podTemplateLister) List(selector labels.Selector) (ret []*v1.PodTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PodTemplate))
	})
	return ret, err
}

// PodTemplates returns an object that can list and get PodTemplates.
func (s *podTemplateLister) PodTemplates(namespace string) PodTemplateNamespaceLister {
	return podTemplateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodTemplateNamespaceLister helps list and get PodTemplates.
type PodTemplateNamespaceLister interface {
	// List lists all PodTemplates in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.PodTemplate, err error)
	// Get retrieves the PodTemplate from the indexer for a given namespace and name.
	Get(name string) (*v1.PodTemplate, error)
}

// podTemplateNamespaceLister implements the PodTemplateNamespaceLister
// interface.
type podTemplateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PodTemplates in the indexer for a given namespace.
func (s podTemplateNamespaceLister) List(selector labels.Selector) (ret []*v1.PodTemplate, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.PodTemplate))
	})
	return ret, err
}

// Get retrieves the PodTemplate from the indexer for a given namespace and name.
func (s podTemplateNamespaceLister) Get(name string) (*v1.PodTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("podtemplate"), name)
	}
	return obj.(*v1.PodTemplate), nil
}

// ReplicationControllerLister helps list ReplicationControllers.
type ReplicationControllerLister interface {
	// List lists all ReplicationControllers in the indexer.
	List(selector labels.Selector) (ret []*v1.ReplicationController, err error)
	// ReplicationControllers returns an object that can list and get ReplicationControllers.
	ReplicationControllers(namespace string) ReplicationControllerNamespaceLister
}

// replicationControllerLister implements the ReplicationControllerLister interface.
type replicationControllerLister struct {
	indexer cache.Indexer
}

// NewReplicationControllerLister returns a new ReplicationControllerLister.
func NewReplicationControllerLister(indexer cache.Indexer) ReplicationControllerLister {
	return &replicationControllerLister{indexer: indexer}
}

// List lists all ReplicationControllers in the indexer.
func (s *replicationControllerLister) List(selector labels.Selector) (ret []*v1.ReplicationController, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ReplicationController))
	})
	return ret, err
}

// ReplicationControllers returns an object that can list and get ReplicationControllers.
func (s *replicationControllerLister) ReplicationControllers(namespace string) ReplicationControllerNamespaceLister {
	return replicationControllerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ReplicationControllerNamespaceLister helps list and get ReplicationControllers.
type ReplicationControllerNamespaceLister interface {
	// List lists all ReplicationControllers in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ReplicationController, err error)
	// Get retrieves the ReplicationController from the indexer for a given namespace and name.
	Get(name string) (*v1.ReplicationController, error)
}

// replicationControllerNamespaceLister implements the ReplicationControllerNamespaceLister
// interface.
type replicationControllerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ReplicationControllers in the indexer for a given namespace.
func (s replicationControllerNamespaceLister) List(selector labels.Selector) (ret []*v1.ReplicationController, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ReplicationController))
	})
	return ret, err
}

// Get retrieves the ReplicationController from the indexer for a given namespace and name.
func (s replicationControllerNamespaceLister) Get(name string) (*v1.ReplicationController, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("replicationcontroller"), name)
	}
	return obj.(*v1.ReplicationController), nil
}

// ResourceQuotaLister helps list ResourceQuotas.
type ResourceQuotaLister interface {
	// List lists all ResourceQuotas in the indexer.
	List(selector labels.Selector) (ret []*v1.ResourceQuota, err error)
	// ResourceQuotas returns an object that can list and get ResourceQuotas.
	ResourceQuotas(namespace string) ResourceQuotaNamespaceLister
}

// resourceQuotaLister implements the ResourceQuotaLister interface.
type resourceQuotaLister struct {
	indexer cache.Indexer
}

// NewResourceQuotaLister returns a new ResourceQuotaLister.
func NewResourceQuotaLister(indexer cache.Indexer) ResourceQuotaLister {
	return &resourceQuotaLister{indexer: indexer}
}

// List lists all ResourceQuotas in the indexer.
func (s *resourceQuotaLister) List(selector labels.Selector) (ret []*v1.ResourceQuota, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ResourceQuota))
	})
	return ret, err
}

// ResourceQuotas returns an object that can list and get ResourceQuotas.
func (s *resourceQuotaLister) ResourceQuotas(namespace string) ResourceQuotaNamespaceLister {
	return resourceQuotaNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ResourceQuotaNamespaceLister helps list and get ResourceQuotas.
type ResourceQuotaNamespaceLister interface {
	// List lists all ResourceQuotas in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ResourceQuota, err error)
	// Get retrieves the ResourceQuota from the indexer for a given namespace and name.
	Get(name string) (*v1.ResourceQuota, error)
}

// resourceQuotaNamespaceLister implements the ResourceQuotaNamespaceLister
// interface.
type resourceQuotaNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ResourceQuotas in the indexer for a given namespace.
func (s resourceQuotaNamespaceLister) List(selector labels.Selector) (ret []*v1.ResourceQuota, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ResourceQuota))
	})
	return ret, err
}

// Get retrieves the ResourceQuota from the indexer for a given namespace and name.
func (s resourceQuotaNamespaceLister) Get(name string) (*v1.ResourceQuota, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("resourcequota"), name)
	}
	return obj.(*v1.ResourceQuota), nil
}

// SecretLister helps list Secrets.
type SecretLister interface {
	// List lists all Secrets in the indexer.
	List(selector labels.Selector) (ret []*v1.Secret, err error)
	// Secrets returns an object that can list and get Secrets.
	Secrets(namespace string) SecretNamespaceLister
}

// secretLister implements the SecretLister interface.
type secretLister struct {
	indexer cache.Indexer
}

// NewSecretLister returns a new SecretLister.
func NewSecretLister(indexer cache.Indexer) SecretLister {
	return &secretLister{indexer: indexer}
}

// List lists all Secrets in the indexer.
func (s *secretLister) List(selector labels.Selector) (ret []*v1.Secret, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Secret))
	})
	return ret, err
}

// Secrets returns an object that can list and get Secrets.
func (s *secretLister) Secrets(namespace string) SecretNamespaceLister {
	return secretNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SecretNamespaceLister helps list and get Secrets.
type SecretNamespaceLister interface {
	// List lists all Secrets in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Secret, err error)
	// Get retrieves the Secret from the indexer for a given namespace and name.
	Get(name string) (*v1.Secret, error)
}

// secretNamespaceLister implements the SecretNamespaceLister
// interface.
type secretNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Secrets in the indexer for a given namespace.
func (s secretNamespaceLister) List(selector labels.Selector) (ret []*v1.Secret, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Secret))
	})
	return ret, err
}

// Get retrieves the Secret from the indexer for a given namespace and name.
func (s secretNamespaceLister) Get(name string) (*v1.Secret, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("secret"), name)
	}
	return obj.(*v1.Secret), nil
}

// ServiceLister helps list Services.
type ServiceLister interface {
	// List lists all Services in the indexer.
	List(selector labels.Selector) (ret []*v1.Service, err error)
	// Services returns an object that can list and get Services.
	Services(namespace string) ServiceNamespaceLister
}

// serviceLister implements the ServiceLister interface.
type serviceLister struct {
	indexer cache.Indexer
}

// NewServiceLister returns a new ServiceLister.
func NewServiceLister(indexer cache.Indexer) ServiceLister {
	return &serviceLister{indexer: indexer}
}

// List lists all Services in the indexer.
func (s *serviceLister) List(selector labels.Selector) (ret []*v1.Service, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Service))
	})
	return ret, err
}

// Services returns an object that can list and get Services.
func (s *serviceLister) Services(namespace string) ServiceNamespaceLister {
	return serviceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceNamespaceLister helps list and get Services.
type ServiceNamespaceLister interface {
	// List lists all Services in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Service, err error)
	// Get retrieves the Service from the indexer for a given namespace and name.
	Get(name string) (*v1.Service, error)
}

// serviceNamespaceLister implements the ServiceNamespaceLister
// interface.
type serviceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Services in the indexer for a given namespace.
func (s serviceNamespaceLister) List(selector labels.Selector) (ret []*v1.Service, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Service))
	})
	return ret, err
}

// Get retrieves the Service from the indexer for a given namespace and name.
func (s serviceNamespaceLister) Get(name string) (*v1.Service, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("service"), name)
	}
	return obj.(*v1.Service), nil
}

// ServiceAccountLister helps list ServiceAccounts.
type ServiceAccountLister interface {
	// List lists all ServiceAccounts in the indexer.
	List(selector labels.Selector) (ret []*v1.ServiceAccount, err error)
	// ServiceAccounts returns an object that can list and get ServiceAccounts.
	ServiceAccounts(namespace string) ServiceAccountNamespaceLister
}

// serviceAccountLister implements the ServiceAccountLister interface.
type serviceAccountLister struct {
	indexer cache.Indexer
}

// NewServiceAccountLister returns a new ServiceAccountLister.
func NewServiceAccountLister(indexer cache.Indexer) ServiceAccountLister {
	return &serviceAccountLister{indexer: indexer}
}

// List lists all ServiceAccounts in the indexer.
func (s *serviceAccountLister) List(selector labels.Selector) (ret []*v1.ServiceAccount, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ServiceAccount))
	})
	return ret, err
}

// ServiceAccounts returns an object that can list and get ServiceAccounts.
func (s *serviceAccountLister) ServiceAccounts(namespace string) ServiceAccountNamespaceLister {
	return serviceAccountNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceAccountNamespaceLister helps list and get ServiceAccounts.
type ServiceAccountNamespaceLister interface {
	// List lists all ServiceAccounts in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.ServiceAccount, err error)
	// Get retrieves the ServiceAccount from the indexer for a given namespace and name.
	Get(name string) (*v1.ServiceAccount, error)
}

// serviceAccountNamespaceLister implements the ServiceAccountNamespaceLister
// interface.
type serviceAccountNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServiceAccounts in the indexer for a given namespace.
func (s serviceAccountNamespaceLister) List(selector labels.Selector) (ret []*v1.ServiceAccount, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ServiceAccount))
	})
	return ret, err
}

// Get retrieves the ServiceAccount from the indexer for a given namespace and name.
func (s serviceAccountNamespaceLister) Get(name string) (*v1.ServiceAccount, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(api.Resource("serviceaccount"), name)
	}
	return obj.(*v1.ServiceAccount), nil
}
