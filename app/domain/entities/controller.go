package entities

import (
	"github.com/jinzhu/copier"
)

type ContainerMap map[string]*Container

type PodRepository interface {
	GetPodsForLabels(labels map[string]string) []*Pod
}

type ResouceValues struct {
	// In Megabytes
	Memory int64
	// In mili CPU
	CPU int64
}

type Resource struct {
	Requests ResouceValues
	Limits   ResouceValues
}

type Pod struct {
	Name string
}

type Container struct {
	Name     string
	Resource Resource
}
type Controller struct {
	Namespace     string
	Name          string
	Containers    ContainerMap
	podRepository PodRepository
}

func NewController0(podRepository PodRepository) *Controller {
	return &Controller{
		podRepository: podRepository,
	}
}

func NewController(
	namespace string,
	name string,
	containers ContainerMap,
	podRepository PodRepository,
) *Controller {
	return &Controller{
		Namespace:     namespace,
		Name:          name,
		Containers:    containers,
		podRepository: podRepository,
	}
}

func (c *Controller) GetPods() []*Pod {
	return c.podRepository.GetPodsForLabels(make(map[string]string))
}

func (c *Controller) FullName() string {
	return c.Namespace + "/" + c.Name
}

func (c *Controller) WithContainers(containers ContainerMap) *Controller {
	var copy Controller
	copier.Copy(&copy, &c)
	copy.Containers = containers

	return &copy
}
