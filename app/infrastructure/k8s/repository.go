package k8s

import (
	"flag"
	"path/filepath"

	logs "github.com/sirupsen/logrus"
	"github.com/totemcaf/k8autotune/app/domain/entities"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type K8SRepository struct {
	appsv1 appsv1.AppsV1Interface
	log    *logs.Entry
}

func NewK8SRepository(home string, log *logs.Entry) *K8SRepository {
	var kubeconfig *string

	log.Infof("Configuring K8S from %s", home)

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

	return &K8SRepository{appsv1: clientset.AppsV1(), log: log}
}

func (r *K8SRepository) FindDeploymentsInNamespace(namespace string) []*entities.Controller {
	opts := metav1.ListOptions{}

	deploymentList, err := r.appsv1.Deployments(namespace).List(opts)

	if err != nil {
		r.log.Error(err)
		return []*entities.Controller{}
	}

	controllers := make([]*entities.Controller, 0, len(deploymentList.Items))

	for _, deployment := range deploymentList.Items {
		controllers = append(controllers, r.deployment2controller(&deployment))
	}

	return controllers
}

func (r *K8SRepository) UpdateResources(controller *entities.Controller) error {
	return nil
}

func (r *K8SRepository) GetPodsForLabels(labels map[string]string) []*entities.Pod {
	return nil
}

func (r *K8SRepository) deployment2controller(d *v1.Deployment) *entities.Controller {
	c := entities.NewController0(r)

	c.Namespace = d.GetNamespace()
	c.Name = d.GetName()

	cs := make(map[string]*entities.Container, len(d.Spec.Template.Spec.Containers))

	for _, container := range d.Spec.Template.Spec.Containers {
		resources := container.Resources
		name := container.Name
		cs[name] = &entities.Container{
			Name: name,
			Resource: entities.Resource{
				Requests: entities.ResouceValues{
					CPU: resources.Requests.Cpu().MilliValue(),
				},
				Limits: entities.ResouceValues{
					CPU: resources.Limits.Cpu().MilliValue(),
				},
			},
		}
	}

	return c
}
