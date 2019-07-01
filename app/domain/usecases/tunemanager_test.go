package usecases

import (
	"context"
	"testing"

	logs "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/mock"
	"github.com/totemcaf/k8autotune/app/domain/entities"
	"github.com/totemcaf/k8autotune/mocks"
)

func TestTuneManager_TuneNamespace_NoDeployments(t *testing.T) {
	// Arrange
	cr := new(mocks.ControllerRepository)
	mr := &mocks.MetricsRepository{}

	tm := NewTuneManager(cr, mr, newLog())

	cr.On("FindDeploymentsInNamespace", namespace1).Return(make([]*entities.Controller, 0, 0))

	// Act
	tm.TuneNamespace(namespace1)

	// Assert
	cr.AssertExpectations(t)
	cr.AssertNumberOfCalls(t, "UpdateResources", 0)
}

func TestTuneManager_TuneNamespace_DeploymentWithNoPods(t *testing.T) {
	// Arrange
	cr := new(mocks.ControllerRepository)
	mr := &mocks.MetricsRepository{}

	tm := NewTuneManager(cr, mr, newLog())

	pr := new(mocks.PodRepository)

	d := entities.NewController0(pr)

	cr.On("FindDeploymentsInNamespace", namespace1).Return([]*entities.Controller{d})
	pr.On("GetPodsForLabels", make(map[string]string)).Return([]*entities.Pod{})

	// Act
	tm.TuneNamespace(namespace1)

	// Assert
	cr.AssertExpectations(t)
	cr.AssertNumberOfCalls(t, "UpdateResources", 0)
}

func TestTuneManager_TuneNamespace_UpdateContainersOfDeployment(t *testing.T) {
	// Arrange
	cr := new(mocks.ControllerRepository)
	mr := new(mocks.MetricsRepository)

	tm := NewTuneManager(cr, mr, newLog())

	pr := new(mocks.PodRepository)

	containers := entities.ContainerMap{
		containerName1: &entities.Container{
			Name: containerName1,
			Resource: entities.Resource{
				Requests: entities.ResouceValues{
					CPU: 10,
				},
				Limits: entities.ResouceValues{
					CPU: 50,
				},
			},
		},
	}
	d := entities.NewController(namespace1, controllerName1, containers, pr)

	cr.On("FindDeploymentsInNamespace", namespace1).Return([]*entities.Controller{d})
	mr.On("CPUUsedForContainer", d, containerName1).Return(makeMetrics(500, 1000))
	cr.On("UpdateResources", mock.AnythingOfType("*entities.Controller")).Return(nil)

	// Act
	tm.TuneNamespace(namespace1)

	// Assert
	cr.AssertExpectations(t)
	cr.AssertNumberOfCalls(t, "UpdateResources", 1)
}

func newLog() *logs.Entry {
	return logs.New().WithContext(context.Background())
}
