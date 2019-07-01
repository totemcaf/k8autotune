package usecases

import (
	"github.com/totemcaf/k8autotune/app/domain/entities"
)

// Tuner .
type Tuner struct {
	minimumValuesToTune int
	maximumCPUMargin    int64
	cpuMargin           int64
}

// NewTuner builds a Tuner
func NewTuner() *Tuner {
	return &Tuner{
		minimumValuesToTune: 672,
		maximumCPUMargin:    200,
		cpuMargin:           20,
	}
}

// Tune computes the new values for the controller based on its use
func (t *Tuner) tune(container *entities.Container, cpuMetrics []int64) *entities.Container {
	if len(cpuMetrics) < t.minimumValuesToTune {
		return container
	}

	cpuRequest := max64(cpuMetrics)
	safetyMargin := cpuRequest * t.cpuMargin / 100
	if safetyMargin > t.maximumCPUMargin {
		cpuRequest += t.maximumCPUMargin
	} else {
		cpuRequest += safetyMargin
	}

	cpuLimit := container.Resource.Limits.CPU

	if cpuLimit <= cpuRequest {
		cpuLimit = cpuRequest * 120 / 100
	}

	if container.Resource.Requests.CPU == cpuRequest &&
		container.Resource.Limits.CPU == cpuLimit {
		// No changes
		return container
	}
	return &entities.Container{
		Name: container.Name,
		Resource: entities.Resource{
			Requests: entities.ResouceValues{CPU: cpuRequest},
			Limits:   entities.ResouceValues{CPU: cpuLimit},
		},
	}
}

func max64(a []int64) int64 {
	if len(a) == 0 {
		return 0
	}

	m := a[0]

	for _, v := range a {
		if m < v {
			m = v
		}
	}

	return m
}
