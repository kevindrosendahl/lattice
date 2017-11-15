package main

import (
	"flag"
	"fmt"
	"time"

	coreconstants "github.com/mlab-lattice/core/pkg/constants"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	localDevDockerRegistry = "lattice-local"
	devDockerRegistry      = "gcr.io/lattice-dev"
)

var (
	kubeconfigPath           string
	provider                 string
	awsRegion                string
	systemIP                 string
	userSystemUrl            string
	latticeContainerRegistry string
	componentBuildRegistry   string
)

func init() {
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "path to kubeconfig file")
	flag.StringVar(&provider, "provider", "", "name of provider to use")
	flag.StringVar(&awsRegion, "aws-region", "", "name of aws region to use")
	flag.StringVar(&systemIP, "system-ip", "", "IP address of the system if -provider=local")
	flag.StringVar(&userSystemUrl, "user-system-url", "", "url of the user-system definition")
	flag.StringVar(&latticeContainerRegistry, "lattice-container-registry", "", "registry used to pull lattice containers")
	flag.StringVar(&componentBuildRegistry, "component-build-registry", "", "registry used to push component builds to")
	flag.Parse()
}

func main() {
	switch provider {
	case coreconstants.ProviderLocal, coreconstants.ProviderAWS:
	default:
		panic("unsupported provider")
	}

	var config *rest.Config
	var err error
	if kubeconfigPath == "" {
		config, err = rest.InClusterConfig()
	} else {
		// TODO: support passing in the context when supported
		// https://github.com/kubernetes/minikube/issues/2100
		//configOverrides := &clientcmd.ConfigOverrides{CurrentContext: kubeContext}
		configOverrides := &clientcmd.ConfigOverrides{}
		config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
			configOverrides,
		).ClientConfig()
	}

	if err != nil {
		panic(err)
	}

	kubeClientset := clientset.NewForConfigOrDie(config)

	seedNamespaces(kubeClientset)
	seedCrds(config)
	seedRbac(kubeClientset)
	seedConfig(config, userSystemUrl, systemIP, awsRegion)
	seedEnvoyXdsApi(kubeClientset)
	seedLatticeControllerManager(kubeClientset)
	seedLatticeSystemEnvironmentManagerAPI(kubeClientset)

	if provider == coreconstants.ProviderLocal {
		seedLocalSpecific(kubeClientset)
	} else {
		seedCloudSpecific(kubeClientset)
	}
}

func pollKubeResourceCreation(resourceCreationFunc func() (interface{}, error)) {
	err := wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		_, err := resourceCreationFunc()

		if err != nil && !apierrors.IsAlreadyExists(err) {
			fmt.Printf("encountered error from API: %v\n", err)
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		panic(err)
	}
}