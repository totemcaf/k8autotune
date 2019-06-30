package usecases

import (
	"context"
	"testing"

	logs "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/mock"
	"github.com/totemcaf/k8autotune.git/app/domain/entities"
	"github.com/totemcaf/k8autotune.git/mocks"
)

const namespace1 = "ns1"

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

	d := entities.NewController(pr)

	cr.On("FindDeploymentsInNamespace", namespace1).Return([]*entities.Controller{d})
	pr.On("GetPodsForLabels", make(map[string]string)).Return([]*entities.Pod{})

	// Act
	tm.TuneNamespace(namespace1)

	// Assert
	cr.AssertExpectations(t)
	cr.AssertNumberOfCalls(t, "UpdateResources", 0)
}

func TestTuneManager_TuneNamespace_UpdatePodOfDeployment(t *testing.T) {
	// Arrange
	cr := new(mocks.ControllerRepository)
	mr := new(mocks.MetricsRepository)

	tm := NewTuneManager(cr, mr, newLog())

	pr := new(mocks.PodRepository)

	d := entities.NewController(pr)

	pod := &entities.Pod{}

	cr.On("FindDeploymentsInNamespace", namespace1).Return([]*entities.Controller{d})
	pr.On("GetPodsForLabels", make(map[string]string)).Return([]*entities.Pod{pod})
	mr.On("CPUUsedForPod", pod).Return(makeMetrics(500, 1000))
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
