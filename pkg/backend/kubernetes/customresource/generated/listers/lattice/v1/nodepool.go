// This file was automatically generated by lister-gen

package v1

import (
	v1 "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NodePoolLister helps list NodePools.
type NodePoolLister interface {
	// List lists all NodePools in the indexer.
	List(selector labels.Selector) (ret []*v1.NodePool, err error)
	// NodePools returns an object that can list and get NodePools.
	NodePools(namespace string) NodePoolNamespaceLister
	NodePoolListerExpansion
}

// nodePoolLister implements the NodePoolLister interface.
type nodePoolLister struct {
	indexer cache.Indexer
}

// NewNodePoolLister returns a new NodePoolLister.
func NewNodePoolLister(indexer cache.Indexer) NodePoolLister {
	return &nodePoolLister{indexer: indexer}
}

// List lists all NodePools in the indexer.
func (s *nodePoolLister) List(selector labels.Selector) (ret []*v1.NodePool, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NodePool))
	})
	return ret, err
}

// NodePools returns an object that can list and get NodePools.
func (s *nodePoolLister) NodePools(namespace string) NodePoolNamespaceLister {
	return nodePoolNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NodePoolNamespaceLister helps list and get NodePools.
type NodePoolNamespaceLister interface {
	// List lists all NodePools in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.NodePool, err error)
	// Get retrieves the NodePool from the indexer for a given namespace and name.
	Get(name string) (*v1.NodePool, error)
	NodePoolNamespaceListerExpansion
}

// nodePoolNamespaceLister implements the NodePoolNamespaceLister
// interface.
type nodePoolNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NodePools in the indexer for a given namespace.
func (s nodePoolNamespaceLister) List(selector labels.Selector) (ret []*v1.NodePool, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.NodePool))
	})
	return ret, err
}

// Get retrieves the NodePool from the indexer for a given namespace and name.
func (s nodePoolNamespaceLister) Get(name string) (*v1.NodePool, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("nodepool"), name)
	}
	return obj.(*v1.NodePool), nil
}
