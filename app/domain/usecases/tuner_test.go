package usecases

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/totemcaf/k8autotune/app/domain/entities"
)

const containerName1 = "cont-one"
const namespace1 = "namespace-one"
const controllerName1 = "controller-one"

func TestTuner_tuneOne_UseMaxPlusMargenIfValueIsSmallerThanMaximunMargen(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := makeContainer(10)

	metrics := makeMetrics(70, tuner.minimumValuesToTune)

	// Act
	containerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.Equal(t, int64(84), containerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_UseMaxPlusMaximumMargenIfValueIsGratherThanMaximunMargen(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := makeContainer(10)

	max := int64(rand.Intn(100) + 1010)
	metrics := makeMetrics(max, tuner.minimumValuesToTune)

	// Act
	containerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.Equal(t, max+200, containerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_NotEnoughMetricsToChange(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	const currentValue int64 = 7
	container := makeContainer(currentValue)
	metrics := makeMetrics(20, tuner.minimumValuesToTune/2)

	// Act
	containerTuned := tuner.tune(container, metrics)

	// Assert
	assert.Equal(t, currentValue, containerTuned.Resource.Requests.CPU)
}

func TestTuner_tuneOne_IncreaseLimit(t *testing.T) {
	// Arrange
	tuner := NewTuner()

	controller := makeContainer(7)
	metrics := makeMetrics(1800, tuner.minimumValuesToTune)

	// Act
	containerTuned := tuner.tune(controller, metrics)

	// Assert
	assert.True(t, containerTuned.Resource.Requests.CPU < containerTuned.Resource.Limits.CPU)
}

func makeInt64(size int, defaultValue int64) []int64 {
	a := make([]int64, size)

	for i := range a {
		a[i] = defaultValue
	}

	return a
}

func makeMetrics(max int64, values int) []int64 {
	return append([]int64{max * 50 / 100, max * 50 / 100, max, max * 90 / 100, max * 80 / 100}, makeInt64(values, 9)...)
}

func makeContainer(currentCPUValue int64) *entities.Container {
	return &entities.Container{
		Resource: entities.Resource{
			Requests: entities.ResouceValues{
				CPU: currentCPUValue,
			},
		},
	}
}
