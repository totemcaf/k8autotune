package usecases

import (
	logs "github.com/sirupsen/logrus"
	"github.com/totemcaf/k8autotune.git/app/domain/entities"
)

type ControllerRepository interface {
	FindDeploymentsInNamespace(namespace string) []*entities.Controller
	UpdateResources(controller *entities.Controller) error
}

type MetricsRepository interface {
	CPUUsedForPod(pod *entities.Pod) []int
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
		tm.tune(controller, log.WithField("controller", "a name"))
	}
}

func (tm *TuneManager) tune(controller *entities.Controller, log *logs.Entry) {
	pods := controller.GetPods()

	if len(pods) > 0 {
		log.WithField("count", len(pods)).Info("Tuning pods")
		metrics := tm.getMetricsForPods(pods)

		tunedController := tm.tuner.tune(controller, metrics[0]) // TODO make multi pod metrics

		if tunedController != controller {
			log.WithFields(logs.Fields{
				"prevCPU": controller.Resource.Requests.CPU,
				"newCPU":  tunedController.Resource.Requests.CPU,
			}).Info("Updating controller")

			tm.controllerRepository.UpdateResources(tunedController)
		}
	} else {
		log.Info("No pods to update")
	}
}

func (tm *TuneManager) getMetricsForPods(pods []*entities.Pod) [][]int {
	metrics := make([][]int, 0, len(pods))

	for _, pod := range pods {
		metrics = append(metrics, tm.metricsRepository.CPUUsedForPod(pod))
	}
	return metrics
}
