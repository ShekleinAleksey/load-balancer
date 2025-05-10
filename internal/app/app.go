package app

import (
	"github.com/ShekleinAleksey/load-balancer/config"
	"github.com/ShekleinAleksey/load-balancer/internal/server"
	"github.com/ShekleinAleksey/load-balancer/pkg/logger"
	"github.com/sirupsen/logrus"
)

func Run() {
	logger.SetLogrus()

	// Загрузка конфига
	config, err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}

	srv, err := server.New(&config)
	if err != nil {
		logrus.Fatalf("Failed to create server: %v", err)
	}

	// Запуск сервера
	if err := srv.Start(); err != nil {
		logrus.Fatalf("Server failed: %v", err)
	}
}
