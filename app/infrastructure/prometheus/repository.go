package prometheus

import "github.com/totemcaf/k8autotune.git/app/domain/entities"

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CPUUsedForPod(pod *entities.Pod) []int {
	return []int{}
}
