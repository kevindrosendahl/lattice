package servicemesh

import (
	"fmt"

	crv1 "github.com/mlab-lattice/system/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	clusterbootstrapper "github.com/mlab-lattice/system/pkg/backend/kubernetes/lifecycle/cluster/bootstrap/bootstrapper"
	systembootstrapper "github.com/mlab-lattice/system/pkg/backend/kubernetes/lifecycle/system/bootstrap/bootstrapper"
	"github.com/mlab-lattice/system/pkg/backend/kubernetes/servicemesh/envoy"

	appsv1 "k8s.io/api/apps/v1"
)

type Interface interface {
	clusterbootstrapper.Interface
	systembootstrapper.Interface

	ServiceAnnotations(*crv1.Service) (map[string]string, error)

	// TransformServiceDeploymentSpec takes in the DeploymentSpec generated for a Service, and applies an service mesh
	// related transforms necessary to a copy of the DeploymentSpec, and returns it.
	TransformServiceDeploymentSpec(*crv1.Service, *appsv1.DeploymentSpec, []*crv1.Service) (*appsv1.DeploymentSpec, error)

	// ServiceMeshPort returns the port the service mesh is listening on for a given component port.
	ServiceMeshPort(*crv1.Service, int32) (int32, error)

	// ServiceMeshPorts returns a map whose keys are component ports and values are the port on which the
	// service mesh is listening on for the given key.
	ServiceMeshPorts(*crv1.Service) (map[int32]int32, error)

	// ServicePort returns the component port for a given port that the service mesh is listening on.
	ServicePort(*crv1.Service, int32) (int32, error)

	// ServiceMeshPorts returns a map whose keys are service mesh ports and values are the component port for
	// which the service mesh is listening on for the given key.
	ServicePorts(*crv1.Service) (map[int32]int32, error)

	// IsDeploymentSpecUpdated checks to see if any part of the current DeploymentSpec that the service mesh is responsible
	// for is out of date compared to the desired deployment spec. If the current DeploymentSpec is current, it also returns
	// a copy of the desired DeploymentSpec with the negation of TransformServiceDeploymentSpec applied.
	// That is, if the aspects of the DeploymentSpec that were transformed by TransformServiceDeploymentSpec are all still
	// current, this method should return true, along with a copy of the DeploymentSpec that should be identical to the
	// DeploymentSpec that was passed in to TransformServiceDeploymentSpec.
	IsDeploymentSpecUpdated(service *crv1.Service, current, desired, untransformed *appsv1.DeploymentSpec) (bool, string, *appsv1.DeploymentSpec)

	GetEndpointSpec(*crv1.ServiceAddress) (*crv1.EndpointSpec, error)
}

func NewServiceMesh(config *crv1.ConfigServiceMesh) (Interface, error) {
	if config.Envoy != nil {
		return envoy.NewEnvoyServiceMesh(config.Envoy), nil
	}

	return nil, fmt.Errorf("no service mesh configuration set")
}
