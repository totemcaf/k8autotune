package k8s

import (
	"flag"
	"path/filepath"

	"github.com/totemcaf/k8autotune.git/app/domain/entities"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SRepository struct {
	appsv1 appsv1.AppsV1Interface
}

func NewK8SRepository(home string) *K8SRepository {
	var kubeconfig *string
	if home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &K8SRepository{appsv1: clientset.AppsV1()}
}

func (r *K8SRepository) FindDeploymentsInNamespace(namespace string) []*entities.Controller {
	return nil
}

func (r *K8SRepository) UpdateResources(controller *entities.Controller) error {
	return nil
}

func (r *K8SRepository) GetPodsForLabels(labels map[string]string) []*entities.Pod {
	return nil
}
