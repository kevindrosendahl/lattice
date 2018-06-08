package v2

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	envoylog "github.com/envoyproxy/go-control-plane/pkg/log"
	envoyserver "github.com/envoyproxy/go-control-plane/pkg/server"

	"github.com/mlab-lattice/lattice/pkg/definition/tree"
)

type Backend interface {
	envoylog.Logger
	envoyserver.Callbacks
	envoycache.NodeHash

	Ready() bool
	Run(threadiness int) error
	XDSCache() envoycache.Cache
	SetXDSCacheSnapshot(id string, endpoints, clusters, routes, listeners []envoycache.Resource) error
	Services(serviceCluster string) (map[tree.NodePath]*Service, error)
}
