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

	controller := makeController(10)

	metrics := makeMetrics(70, tuner.minimumValuesToTune)

	// Act
	controllerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.Equal(t, 84, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_UseMaxPlusMaximumMargenIfValueIsGratherThanMaximunMargen(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := makeController(10)

	max := rand.Intn(100) + 1010
	metrics := makeMetrics(max, tuner.minimumValuesToTune)

	// Act
	controllerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.Equal(t, max+200, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_NotEnoughMetricsToChange(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	currentValue := 7
	controller := makeController(currentValue)
	metrics := makeMetrics(20, tuner.minimumValuesToTune/2)

	// Act
	controllerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.Equal(t, currentValue, controllerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_IncreaseLimit(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := makeController(7)
	metrics := makeMetrics(1800, tuner.minimumValuesToTune)

	// Act
	controllerTuned := tuner.tune(controller, metrics)

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

func makeMetrics(max int, values int) []int {
	return append([]int{max * 50 / 100, max * 50 / 100, max, max * 90 / 100, max * 80 / 100}, makeInt(values, 9)...)
}

func makeController(currentCPUValue int) *entities.Controller {
	return &entities.Controller{
		Resource: entities.Resource{
			Requests: entities.ResouceValues{
				CPU: currentCPUValue,
			},
		},
	}
}
