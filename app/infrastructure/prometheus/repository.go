package prometheus

import "github.com/totemcaf/k8autotune/app/domain/entities"

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CPUUsedForContainer(container *entities.Controller, containerName string) []int64 {
	return []int64{}
}
