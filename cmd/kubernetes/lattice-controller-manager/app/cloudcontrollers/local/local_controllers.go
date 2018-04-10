package aws

import (
	controller "github.com/mlab-lattice/lattice/cmd/kubernetes/lattice-controller-manager/app/common"
)

func GetControllerInitializers() map[string]controller.Initializer {
	return map[string]controller.Initializer{
		"load-balancer": initializeLoadBalancerController,
		"node-pool":     initializeNodePoolController,
	}
}
