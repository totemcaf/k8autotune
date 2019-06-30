package entities

type PodRepository interface {
	GetPodsForLabels(labels map[string]string) []*Pod
}

type ResouceValues struct {
	// In Megabytes
	Memory int64
	// In mili CPU
	CPU int
}

type Resource struct {
	Requests ResouceValues
	Limits   ResouceValues
}

type Pod struct {
	Name string
}

type Controller struct {
	Resource      Resource
	podRepository PodRepository
}

func NewController(podRepository PodRepository) *Controller {
	return &Controller{podRepository: podRepository}
}

func (c *Controller) GetPods() []*Pod {
	return c.podRepository.GetPodsForLabels(make(map[string]string))
}
