package usecases

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/totemcaf/k8autotune.git/app/domain/entities"
)

func TestTuner_tuneOne_UseMaxPlusMargenIfValueIsSmallerThanMaximunMargen(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := &entities.Controller{}
	max := 70
	metrics := append([]int{5, 5, max, 20, 10}, makeInt(tuner.minimumValuesToTune, 9)...)

	// Act
	controllerTuned := tuner.Tune(controller, metrics)

	// Assert
	assert.Equal(t, 84, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_UseMaxPlusMaximumMargenIfValueIsGratherThanMaximunMargen(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := &entities.Controller{}
	max := rand.Intn(100) + 1010
	metrics := append([]int{500, 500, max, 200, 100}, makeInt(tuner.minimumValuesToTune, 9)...)

	// Act
	controllerTuned := tuner.Tune(controller, metrics)

	// Assert
	assert.Equal(t, max+200, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_NotEnoughMetricsToChange(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	currentValue := 7
	controller := &entities.Controller{Resource: entities.Resource{Requests: entities.ResouceValues{CPU: currentValue}}}
	metrics := []int{5, 5, 8, 20, 10}

	// Act
	controllerTuned := tuner.Tune(controller, metrics)

	// Assert
	assert.Equal(t, currentValue, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_IncreaseLimit(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	currentValue := 7
	controller := &entities.Controller{Resource: entities.Resource{Requests: entities.ResouceValues{CPU: currentValue}}}
	metrics := append([]int{500, 500, 1800, 200, 100}, makeInt(tuner.minimumValuesToTune, 9)...)

	// Act
	controllerTuned := tuner.Tune(controller, metrics)

	// Assert
	assert.True(t, controllerTuned.Resource.Requests.CPU < controllerTuned.Resource.Limits.CPU)
}

func makeInt(size int, defaultValue int) []int {
	a := make([]int, size)

	for i := range a {
		a[i] = defaultValue
	}

	return a
}
