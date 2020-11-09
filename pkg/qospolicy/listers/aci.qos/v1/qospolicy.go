/***
Copyright 2019 Cisco Systems Inc. All rights reserved.

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

package v1

import (
	v1 "github.com/noironetworks/aci-containers/pkg/qospolicy/apis/aci.qos/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// QosPolicyLister helps list QosPolicies.
type QosPolicyLister interface {
	// List lists all QosPolicies in the indexer.
	List(selector labels.Selector) (ret []*v1.QosPolicy, err error)
	// Get retrieves the QosPolicy from the index for a given name.
	Get(name string) (*v1.QosPolicy, error)
	QosPolicyListerExpansion
}

// qosPolicyLister implements the QosPolicyLister interface.
type qosPolicyLister struct {
	indexer cache.Indexer
}

// NewQosPolicyLister returns a new QosPolicyLister.
func NewQosPolicyLister(indexer cache.Indexer) QosPolicyLister {
	return &qosPolicyLister{indexer: indexer}
}

// List lists all QosPolicies in the indexer.
func (s *qosPolicyLister) List(selector labels.Selector) (ret []*v1.QosPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.QosPolicy))
	})
	return ret, err
}

// Get retrieves the QosPolicy from the index for a given name.
func (s *qosPolicyLister) Get(name string) (*v1.QosPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("qospolicy"), name)
	}
	return obj.(*v1.QosPolicy), nil
}
