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

package internalversion

import (
	api "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	rbac "k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/labels"
)

// ClusterRoleLister helps list ClusterRoles.
type ClusterRoleLister interface {
	// List lists all ClusterRoles in the indexer.
	List(selector labels.Selector) (ret []*rbac.ClusterRole, err error)
	// Get retrieves the ClusterRole from the index for a given name.
	Get(name string) (*rbac.ClusterRole, error)
}

// clusterRoleLister implements the ClusterRoleLister interface.
type clusterRoleLister struct {
	indexer cache.Indexer
}

// NewClusterRoleLister returns a new ClusterRoleLister.
func NewClusterRoleLister(indexer cache.Indexer) ClusterRoleLister {
	return &clusterRoleLister{indexer: indexer}
}

// List lists all ClusterRoles in the indexer.
func (s *clusterRoleLister) List(selector labels.Selector) (ret []*rbac.ClusterRole, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.ClusterRole))
	})
	return ret, err
}

// Get retrieves the ClusterRole from the index for a given name.
func (s *clusterRoleLister) Get(name string) (*rbac.ClusterRole, error) {
	key := &rbac.ClusterRole{ObjectMeta: api.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(rbac.Resource("clusterrole"), name)
	}
	return obj.(*rbac.ClusterRole), nil
}

// ClusterRoleBindingLister helps list ClusterRoleBindings.
type ClusterRoleBindingLister interface {
	// List lists all ClusterRoleBindings in the indexer.
	List(selector labels.Selector) (ret []*rbac.ClusterRoleBinding, err error)
	// Get retrieves the ClusterRoleBinding from the index for a given name.
	Get(name string) (*rbac.ClusterRoleBinding, error)
}

// clusterRoleBindingLister implements the ClusterRoleBindingLister interface.
type clusterRoleBindingLister struct {
	indexer cache.Indexer
}

// NewClusterRoleBindingLister returns a new ClusterRoleBindingLister.
func NewClusterRoleBindingLister(indexer cache.Indexer) ClusterRoleBindingLister {
	return &clusterRoleBindingLister{indexer: indexer}
}

// List lists all ClusterRoleBindings in the indexer.
func (s *clusterRoleBindingLister) List(selector labels.Selector) (ret []*rbac.ClusterRoleBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.ClusterRoleBinding))
	})
	return ret, err
}

// Get retrieves the ClusterRoleBinding from the index for a given name.
func (s *clusterRoleBindingLister) Get(name string) (*rbac.ClusterRoleBinding, error) {
	key := &rbac.ClusterRoleBinding{ObjectMeta: api.ObjectMeta{Name: name}}
	obj, exists, err := s.indexer.Get(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(rbac.Resource("clusterrolebinding"), name)
	}
	return obj.(*rbac.ClusterRoleBinding), nil
}

// RoleLister helps list Roles.
type RoleLister interface {
	// List lists all Roles in the indexer.
	List(selector labels.Selector) (ret []*rbac.Role, err error)
	// Roles returns an object that can list and get Roles.
	Roles(namespace string) RoleNamespaceLister
}

// roleLister implements the RoleLister interface.
type roleLister struct {
	indexer cache.Indexer
}

// NewRoleLister returns a new RoleLister.
func NewRoleLister(indexer cache.Indexer) RoleLister {
	return &roleLister{indexer: indexer}
}

// List lists all Roles in the indexer.
func (s *roleLister) List(selector labels.Selector) (ret []*rbac.Role, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.Role))
	})
	return ret, err
}

// Roles returns an object that can list and get Roles.
func (s *roleLister) Roles(namespace string) RoleNamespaceLister {
	return roleNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RoleNamespaceLister helps list and get Roles.
type RoleNamespaceLister interface {
	// List lists all Roles in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*rbac.Role, err error)
	// Get retrieves the Role from the indexer for a given namespace and name.
	Get(name string) (*rbac.Role, error)
}

// roleNamespaceLister implements the RoleNamespaceLister
// interface.
type roleNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Roles in the indexer for a given namespace.
func (s roleNamespaceLister) List(selector labels.Selector) (ret []*rbac.Role, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.Role))
	})
	return ret, err
}

// Get retrieves the Role from the indexer for a given namespace and name.
func (s roleNamespaceLister) Get(name string) (*rbac.Role, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(rbac.Resource("role"), name)
	}
	return obj.(*rbac.Role), nil
}

// RoleBindingLister helps list RoleBindings.
type RoleBindingLister interface {
	// List lists all RoleBindings in the indexer.
	List(selector labels.Selector) (ret []*rbac.RoleBinding, err error)
	// RoleBindings returns an object that can list and get RoleBindings.
	RoleBindings(namespace string) RoleBindingNamespaceLister
}

// roleBindingLister implements the RoleBindingLister interface.
type roleBindingLister struct {
	indexer cache.Indexer
}

// NewRoleBindingLister returns a new RoleBindingLister.
func NewRoleBindingLister(indexer cache.Indexer) RoleBindingLister {
	return &roleBindingLister{indexer: indexer}
}

// List lists all RoleBindings in the indexer.
func (s *roleBindingLister) List(selector labels.Selector) (ret []*rbac.RoleBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.RoleBinding))
	})
	return ret, err
}

// RoleBindings returns an object that can list and get RoleBindings.
func (s *roleBindingLister) RoleBindings(namespace string) RoleBindingNamespaceLister {
	return roleBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// RoleBindingNamespaceLister helps list and get RoleBindings.
type RoleBindingNamespaceLister interface {
	// List lists all RoleBindings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*rbac.RoleBinding, err error)
	// Get retrieves the RoleBinding from the indexer for a given namespace and name.
	Get(name string) (*rbac.RoleBinding, error)
}

// roleBindingNamespaceLister implements the RoleBindingNamespaceLister
// interface.
type roleBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all RoleBindings in the indexer for a given namespace.
func (s roleBindingNamespaceLister) List(selector labels.Selector) (ret []*rbac.RoleBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*rbac.RoleBinding))
	})
	return ret, err
}

// Get retrieves the RoleBinding from the indexer for a given namespace and name.
func (s roleBindingNamespaceLister) Get(name string) (*rbac.RoleBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(rbac.Resource("rolebinding"), name)
	}
	return obj.(*rbac.RoleBinding), nil
}
