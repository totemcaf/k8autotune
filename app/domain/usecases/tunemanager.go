package usecases

import (
	logs "github.com/sirupsen/logrus"
	"github.com/totemcaf/k8autotune/app/domain/entities"
)

type ControllerRepository interface {
	FindDeploymentsInNamespace(namespace string) []*entities.Controller
	UpdateResources(controller *entities.Controller) error
}

type MetricsRepository interface {
	CPUUsedForContainer(container *entities.Controller, containerName string) []int64 // NOooooo
}

type TuneManager struct {
	controllerRepository ControllerRepository
	metricsRepository    MetricsRepository
	tuner                *Tuner
	log                  *logs.Entry
}

func NewTuneManager(controllerRepository ControllerRepository, metricsRepository MetricsRepository, log *logs.Entry) *TuneManager {
	return &TuneManager{controllerRepository: controllerRepository, metricsRepository: metricsRepository, tuner: NewTuner(), log: log}
}

func (tm *TuneManager) TuneNamespace(namespace string) {
	log := tm.log.WithField("namespace", namespace)
	log.Info("Start tunning")

	controllers := tm.controllerRepository.FindDeploymentsInNamespace(namespace)

	log.WithField("count", len(controllers)).Info("Found controllers")
	for _, controller := range controllers {
		tm.tune(controller, log.WithField("controller", controller.FullName()))
	}
}

func (tm *TuneManager) tune(controller *entities.Controller, log *logs.Entry) {

	log.WithField("count", len(controller.Containers)).Info("Tuning containers")

	changed := false
	tunnedContainers := make(entities.ContainerMap, len(controller.Containers))

	for _, container := range controller.Containers {
		metrics := tm.metricsRepository.CPUUsedForContainer(controller, container.Name)

		tunedContainer := tm.tuner.tune(container, metrics) // TODO make multi pod metrics

		changed = changed || tunedContainer != container

		tunnedContainers[tunedContainer.Name] = tunedContainer
	}

	if changed {
		tunedController := controller.WithContainers(tunnedContainers)
		for _, c := range tunnedContainers {
			log.WithFields(logs.Fields{
				"container": c.Name,
				"prevCPU":   controller.Containers[c.Name].Resource.Requests.CPU,
				"newCPU":    c.Resource.Requests.CPU,
			}).Info("Updating container in controller")
		}

		tm.controllerRepository.UpdateResources(tunedController)
	} else {
		log.Info("No pods to update")
	}
}
