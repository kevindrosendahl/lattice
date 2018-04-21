// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=lattice.mlab.com, Version=v1
	case v1.SchemeGroupVersion.WithResource("addresses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Addresses().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("builds"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Builds().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("componentbuilds"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().ComponentBuilds().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("configs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Configs().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("deploys"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Deploys().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("loadbalancers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().LoadBalancers().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("nodepools"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().NodePools().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("services"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Services().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("servicebuilds"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().ServiceBuilds().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("systems"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Systems().Informer()}, nil
	case v1.SchemeGroupVersion.WithResource("teardowns"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Lattice().V1().Teardowns().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
