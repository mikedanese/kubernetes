/*
Copyright 2014 The Kubernetes Authors All rights reserved.
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
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/validation"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/fields"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/registry/generic"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/runtime"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util/fielderrors"
)

// jobStrategy implements verification logic for Replication Controllers.
type jobStrategy struct {
	runtime.ObjectTyper
	api.NameGenerator
}

// Strategy is the default logic that applies when creating and updating Replication Controller objects.
var Strategy = jobStrategy{api.Scheme, api.SimpleNameGenerator}

// NamespaceScoped returns true because all Replication Controllers need to be within a namespace.
func (jobStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate clears the status of a replication controller before creation.
func (jobStrategy) PrepareForCreate(obj runtime.Object) {
	controller := obj.(*api.ReplicationController)
	controller.Status = api.ReplicationControllerStatus{}

	controller.Generation = 1
	controller.Status.ObservedGeneration = 0
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
func (jobStrategy) PrepareForUpdate(obj, old runtime.Object) {
	// TODO: once JOB has a status sub-resource we can enable this.
	//newController := obj.(*api.ReplicationController)
	//oldController := old.(*api.ReplicationController)
	//newController.Status = oldController.Status
	newController := obj.(*api.Job)
	oldController := old.(*api.Job)

	// Any changes to the spec increment the generation number, any changes to the
	// status should reflect the generation number of the corresponding object. We push
	// the burden of managing the status onto the clients because we can't (in general)
	// know here what version of spec the writer of the status has seen. It may seem like
	// we can at first -- since obj contains spec -- but in the future we will probably make
	// status its own object, and even if we don't, writes may be the result of a
	// read-update-write loop, so the contents of spec may not actually be the spec that
	// the controller has *seen*.
	//
	// TODO: Any changes to a part of the object that represents desired state (labels,
	// annotations etc) should also increment the generation.
	if !reflect.DeepEqual(oldController.Spec, newController.Spec) {
		newController.Generation = oldController.Generation + 1
	}
}

// Validate validates a new replication controller.
func (jobStrategy) Validate(ctx api.Context, obj runtime.Object) fielderrors.ValidationErrorList {
	controller := obj.(*api.Job)
	return validation.ValidateJob(controller)
}

// AllowCreateOnUpdate is false for replication controllers; this means a POST is
// needed to create one.
func (jobStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
func (jobStrategy) ValidateUpdate(ctx api.Context, obj, old runtime.Object) fielderrors.ValidationErrorList {
	validationErrorList := validation.ValidateJob(obj.(*api.Job))
	updateErrorList := validation.ValidateJobUpdate(old.(*api.Job), obj.(*api.Job))
	return append(validationErrorList, updateErrorList...)
}

func (jobStrategy) AllowUnconditionalUpdate() bool {
	return true
}

// MatchJob returns a generic matcher for a given label and field selector.
func MatchJob(label labels.Selector, field fields.Selector) generic.Matcher {
	return &generic.SelectionPredicate{
		Label: label,
		Field: field,
		GetAttrs: func(obj runtime.Object) (labels.Set, fields.Set, error) {
			job, ok := obj.(*api.Job)
			if !ok {
				return nil, nil, fmt.Errorf("not a job")
			}
			return labels.Set(job.ObjectMeta.Labels), JobToSelectableFields(job), nil
		},
	}
}

func JobToSelectableFields(job *api.Job) fields.Set {
	return fields.Set{
		"metadata.name": job.Name,
	}
}
