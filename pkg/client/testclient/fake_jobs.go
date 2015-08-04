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

package testclient

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/watch"
)

// FakeJobs implements JobsInterface. Meant to be embedded into a struct to get a default
// implementation. This makes faking out just the methods you want to test easier.
type FakeJobs struct {
	Fake      *Fake
	Namespace string
}

func (c *FakeJobs) List(label labels.Selector) (*api.JobList, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "list-jobs"}, &api.JobList{})
	return obj.(*api.JobList), err
}

func (c *FakeJobs) Get(name string) (*api.Job, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "get-job", Value: name}, &api.Job{})
	return obj.(*api.Job), err
}

func (c *FakeJobs) Delete(name string) error {
	_, err := c.Fake.Invokes(FakeAction{Action: "delete-job", Value: name}, &api.Job{})
	return err
}

func (c *FakeJobs) Create(job *api.Job) (*api.Job, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "create-job"}, &api.Job{})
	return obj.(*api.Job), err
}

func (c *FakeJobs) Update(job *api.Job) (*api.Job, error) {
	obj, err := c.Fake.Invokes(FakeAction{Action: "update-job", Value: job.Name}, &api.Job{})
	return obj.(*api.Job), err
}

func (c *FakeJobs) Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error) {
	c.Fake.Invokes(FakeAction{Action: "watch-jobs", Value: resourceVersion}, nil)
	return c.Fake.Watch, c.Fake.Err()
}
