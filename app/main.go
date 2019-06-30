package main

import (
	"context"
	"os"

	"github.com/totemcaf/k8autotune.git/app/domain/usecases"
	"github.com/totemcaf/k8autotune.git/app/infrastructure/k8s"
	"github.com/totemcaf/k8autotune.git/app/infrastructure/prometheus"

	logs "github.com/sirupsen/logrus"
)

func main() {
	log := logs.New().WithContext(context.Background())

	log.Infof("Starting application v%s", version)

	repository := k8s.NewK8SRepository("")

	metricsRepository := prometheus.NewRepository()

	tuneManager := usecases.NewTuneManager(repository, metricsRepository, log)

	tuneManager.TuneNamespace("default")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
