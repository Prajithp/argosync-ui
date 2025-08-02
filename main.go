package main

import (
	"os"

	"github.com/Prajithp/argosync/internal/api"
	"github.com/Prajithp/argosync/internal/config"
	"github.com/Prajithp/argosync/internal/logger"
	"github.com/Prajithp/argosync/internal/repository"
	"github.com/Prajithp/argosync/internal/server"
	"github.com/Prajithp/argosync/internal/service"
)

//go:generate go run ./hack/gen/main.go

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize logger
	logger.InitLogger(cfg.Debug)
	logger.Info().Msg("Starting Heirloom Deployment API")

	// Initialize repository
	repo, err := repository.InitRepositoryFromConfig(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize repository")
		os.Exit(1)
	}

	defer repo.Close()

	// Create service
	deploymentService := service.NewDeploymentService(repo)

	// Create API handler
	handler := api.NewHandler(deploymentService)

	// Create and start server
	srv := server.NewServer(cfg, handler)
	srv.Start()
}
