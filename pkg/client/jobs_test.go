/*
Copyright 2015 Google Inc. All rights reserved.

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

package client

import (
	"testing"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/testapi"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
)

func TestListJobs(t *testing.T) {
	ns := api.NamespaceAll
	c := &testClient{
		Request: testRequest{
			Method: "GET",
			Path:   testapi.ResourcePath("jobs", ns, ""),
		},
		Response: Response{StatusCode: 200,
			Body: &api.JobList{
				Items: []api.Job{
					{
						ObjectMeta: api.ObjectMeta{
							Name: "foo",
							Labels: map[string]string{
								"foo":  "bar",
								"name": "baz",
							},
						},
						Spec: api.JobSpec{
							Completions: 2,
							Template:    &api.PodTemplateSpec{},
						},
					},
				},
			},
		},
	}
	receivedJobList, err := c.Setup().Jobs(ns).List(labels.Everything())
	c.Validate(t, receivedJobList, err)

}

func TestGetJob(t *testing.T) {
	ns := api.NamespaceDefault
	c := &testClient{
		Request: testRequest{Method: "GET", Path: testapi.ResourcePath("jobs", ns, "foo"), Query: buildQueryValues(ns, nil)},
		Response: Response{
			StatusCode: 200,
			Body: &api.Job{
				ObjectMeta: api.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"foo":  "bar",
						"name": "baz",
					},
				},
				Spec: api.JobSpec{
					Completions: 2,
					Template:    &api.PodTemplateSpec{},
				},
			},
		},
	}
	receivedJob, err := c.Setup().Jobs(ns).Get("foo")
	c.Validate(t, receivedJob, err)
}

func TestGetJobWithNoName(t *testing.T) {
	ns := api.NamespaceDefault
	c := &testClient{Error: true}
	receivedPod, err := c.Setup().Jobs(ns).Get("")
	if (err != nil) && (err.Error() != nameRequiredError) {
		t.Errorf("Expected error: %v, but got %v", nameRequiredError, err)
	}

	c.Validate(t, receivedPod, err)
}

func TestUpdateJob(t *testing.T) {
	ns := api.NamespaceDefault
	requestJob := &api.Job{
		ObjectMeta: api.ObjectMeta{Name: "foo", ResourceVersion: "1"},
	}
	c := &testClient{
		Request: testRequest{Method: "PUT", Path: testapi.ResourcePath("jobs", ns, "foo"), Query: buildQueryValues(ns, nil)},
		Response: Response{
			StatusCode: 200,
			Body: &api.Job{
				ObjectMeta: api.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"foo":  "bar",
						"name": "baz",
					},
				},
				Spec: api.JobSpec{
					Completions: 2,
					Template:    &api.PodTemplateSpec{},
				},
			},
		},
	}
	receivedJob, err := c.Setup().Jobs(ns).Update(requestJob)
	c.Validate(t, receivedJob, err)
}

func TestDeleteJob(t *testing.T) {
	ns := api.NamespaceDefault
	c := &testClient{
		Request:  testRequest{Method: "DELETE", Path: testapi.ResourcePath("jobs", ns, "foo"), Query: buildQueryValues(ns, nil)},
		Response: Response{StatusCode: 200},
	}
	err := c.Setup().Jobs(ns).Delete("foo")
	c.Validate(t, nil, err)
}

func TestCreateJob(t *testing.T) {
	ns := api.NamespaceDefault
	requestJob := &api.Job{
		ObjectMeta: api.ObjectMeta{Name: "foo"},
	}
	c := &testClient{
		Request: testRequest{Method: "POST", Path: testapi.ResourcePath("jobs", ns, ""), Body: requestJob, Query: buildQueryValues(ns, nil)},
		Response: Response{
			StatusCode: 200,
			Body: &api.Job{
				ObjectMeta: api.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"foo":  "bar",
						"name": "baz",
					},
				},
				Spec: api.JobSpec{
					Completions: 2,
					Template:    &api.PodTemplateSpec{},
				},
			},
		},
	}
	receivedJob, err := c.Setup().Jobs(ns).Create(requestJob)
	c.Validate(t, receivedJob, err)
}
