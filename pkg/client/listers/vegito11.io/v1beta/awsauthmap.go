/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta

import (
	v1beta "github.com/vegito11/AWSAuthSync/pkg/apis/vegito11.io/v1beta"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AWSAuthMapLister helps list AWSAuthMaps.
// All objects returned here must be treated as read-only.
type AWSAuthMapLister interface {
	// List lists all AWSAuthMaps in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta.AWSAuthMap, err error)
	// AWSAuthMaps returns an object that can list and get AWSAuthMaps.
	AWSAuthMaps(namespace string) AWSAuthMapNamespaceLister
	AWSAuthMapListerExpansion
}

// aWSAuthMapLister implements the AWSAuthMapLister interface.
type aWSAuthMapLister struct {
	indexer cache.Indexer
}

// NewAWSAuthMapLister returns a new AWSAuthMapLister.
func NewAWSAuthMapLister(indexer cache.Indexer) AWSAuthMapLister {
	return &aWSAuthMapLister{indexer: indexer}
}

// List lists all AWSAuthMaps in the indexer.
func (s *aWSAuthMapLister) List(selector labels.Selector) (ret []*v1beta.AWSAuthMap, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta.AWSAuthMap))
	})
	return ret, err
}

// AWSAuthMaps returns an object that can list and get AWSAuthMaps.
func (s *aWSAuthMapLister) AWSAuthMaps(namespace string) AWSAuthMapNamespaceLister {
	return aWSAuthMapNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// AWSAuthMapNamespaceLister helps list and get AWSAuthMaps.
// All objects returned here must be treated as read-only.
type AWSAuthMapNamespaceLister interface {
	// List lists all AWSAuthMaps in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta.AWSAuthMap, err error)
	// Get retrieves the AWSAuthMap from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta.AWSAuthMap, error)
	AWSAuthMapNamespaceListerExpansion
}

// aWSAuthMapNamespaceLister implements the AWSAuthMapNamespaceLister
// interface.
type aWSAuthMapNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all AWSAuthMaps in the indexer for a given namespace.
func (s aWSAuthMapNamespaceLister) List(selector labels.Selector) (ret []*v1beta.AWSAuthMap, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta.AWSAuthMap))
	})
	return ret, err
}

// Get retrieves the AWSAuthMap from the indexer for a given namespace and name.
func (s aWSAuthMapNamespaceLister) Get(name string) (*v1beta.AWSAuthMap, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta.Resource("awsauthmap"), name)
	}
	return obj.(*v1beta.AWSAuthMap), nil
}
