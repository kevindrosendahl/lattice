package backend

import (
	kubeconstants "github.com/mlab-lattice/system/pkg/backend/kubernetes/constants"
	crv1 "github.com/mlab-lattice/system/pkg/backend/kubernetes/customresource/apis/lattice/v1"
	kubeutil "github.com/mlab-lattice/system/pkg/backend/kubernetes/util/kubernetes"
	"github.com/mlab-lattice/system/pkg/types"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/satori/go.uuid"
)

func (kb *KubernetesBackend) TearDownSystem(id types.SystemID) (types.SystemTeardownID, error) {
	systemTeardown, err := getSystemTeardown(id)
	if err != nil {
		return "", err
	}

	namespace := kubeutil.SystemNamespace(kb.ClusterID, id)
	result, err := kb.LatticeClient.LatticeV1().SystemTeardowns(namespace).Create(systemTeardown)
	if err != nil {
		return "", err
	}

	return types.SystemTeardownID(result.Name), err
}

func getSystemTeardown(id types.SystemID) (*crv1.SystemTeardown, error) {
	labels := map[string]string{
		kubeconstants.LatticeNamespaceLabel: string(id),
	}

	sysT := &crv1.SystemTeardown{
		ObjectMeta: metav1.ObjectMeta{
			Name:   uuid.NewV4().String(),
			Labels: labels,
		},
		Spec: crv1.SystemTeardownSpec{},
		Status: crv1.SystemTeardownStatus{
			State: crv1.SystemTeardownStatePending,
		},
	}

	return sysT, nil
}

func (kb *KubernetesBackend) GetSystemTeardown(id types.SystemID, tid types.SystemTeardownID) (*types.SystemTeardown, bool, error) {
	namespace := kubeutil.SystemNamespace(kb.ClusterID, id)
	result, err := kb.LatticeClient.LatticeV1().SystemTeardowns(namespace).Get(string(tid), metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	sb := &types.SystemTeardown{
		ID:    tid,
		State: getSystemTeardownState(result.Status.State),
	}

	return sb, true, nil
}

func (kb *KubernetesBackend) ListSystemTeardowns(id types.SystemID) ([]types.SystemTeardown, error) {
	namespace := kubeutil.SystemNamespace(kb.ClusterID, id)
	result, err := kb.LatticeClient.LatticeV1().SystemTeardowns(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var teardowns []types.SystemTeardown
	for _, b := range result.Items {
		teardowns = append(teardowns, types.SystemTeardown{
			ID:    types.SystemTeardownID(b.Name),
			State: getSystemTeardownState(b.Status.State),
		})
	}

	return teardowns, nil
}

func getSystemTeardownState(state crv1.SystemTeardownState) types.SystemTeardownState {
	switch state {
	case crv1.SystemTeardownStatePending:
		return types.SystemTeardownStatePending
	case crv1.SystemTeardownStateInProgress:
		return types.SystemTeardownStateInProgress
	case crv1.SystemTeardownStateSucceeded:
		return types.SystemTeardownStateSucceeded
	case crv1.SystemTeardownStateFailed:
		return types.SystemTeardownStateFailed
	default:
		panic("unreachable")
	}
}