// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	lattice_v1 "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	versioned "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/generated/clientset/versioned"
	internalinterfaces "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/generated/informers/externalversions/internalinterfaces"
	v1 "github.com/mlab-lattice/lattice/pkg/backend/kubernetes/customresource/generated/listers/lattice/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ComponentBuildInformer provides access to a shared informer and lister for
// ComponentBuilds.
type ComponentBuildInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ComponentBuildLister
}

type componentBuildInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewComponentBuildInformer constructs a new informer for ComponentBuild type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewComponentBuildInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredComponentBuildInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredComponentBuildInformer constructs a new informer for ComponentBuild type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredComponentBuildInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LatticeV1().ComponentBuilds(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.LatticeV1().ComponentBuilds(namespace).Watch(options)
			},
		},
		&lattice_v1.ComponentBuild{},
		resyncPeriod,
		indexers,
	)
}

func (f *componentBuildInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredComponentBuildInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *componentBuildInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&lattice_v1.ComponentBuild{}, f.defaultInformer)
}

func (f *componentBuildInformer) Lister() v1.ComponentBuildLister {
	return v1.NewComponentBuildLister(f.Informer().GetIndexer())
}
