package main

import (
	"task/configs"
	httpserver "task/internal/adapter/http"
	"task/internal/domain/repositories"
	logger "task/internal/infra"
	"task/internal/usecase"
)

func main() {
	configFilePath := "../configs/config.yaml"
	config, err := configs.LoadConfig(configFilePath)
	if err != nil {
		panic(err)
	}

	log, err := logger.SetupLog(config.Logger.LogFile)
	if err != nil {
		log.Fatal("failed to setup logger", err)
	}

	repo, err := repositories.NewConnectionRepository(config, log)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	service := usecase.NewConnectionService(repo, log)
	server := httpserver.NewServer(service, log)

	log.Info("Server starting on port 8080")
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Recovered from panic: %v", err)
		}
	}()

	if err := server.Start(":8080"); err != nil {
		log.WithError(err).Fatal("Failed to start the server")
	}
}
