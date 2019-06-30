package usecases

import (
	"github.com/totemcaf/k8autotune.git/app/domain/entities"
)

// Tuner .
type Tuner struct {
	minimumValuesToTune int
	maximumCPUMargin    int
	cpuMargin           int
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
func (t *Tuner) tune(controller *entities.Controller, cpuMetrics []int) *entities.Controller {
	if len(cpuMetrics) < t.minimumValuesToTune {
		return controller
	}

	cpuRequest := max(cpuMetrics)
	safetyMargin := cpuRequest * t.cpuMargin / 100
	if safetyMargin > t.maximumCPUMargin {
		cpuRequest += t.maximumCPUMargin
	} else {
		cpuRequest += safetyMargin
	}

	cpuLimit := controller.Resource.Limits.CPU

	if cpuLimit <= cpuRequest {
		cpuLimit = cpuRequest * 120 / 100
	}

	return &entities.Controller{
		Resource: entities.Resource{
			Requests: entities.ResouceValues{CPU: cpuRequest},
			Limits:   entities.ResouceValues{CPU: cpuLimit},
		},
	}
}

func max(a []int) int {
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
