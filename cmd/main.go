package main

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/logger"
	"rinha-de-backend-2025/internal/router"
	"rinha-de-backend-2025/internal/workers"
)

func main() {
	logger := logger.GetLogger("[MAIN]")

	logger.Info("Starting project")

	if err := db.StartDB(); err != nil {
		logger.Fatalf("Failed to start Redis: %v", err)
	}

	if err := db.StartPG(); err != nil {
		logger.Fatalf("Failed to start PostgreSQL: %v", err)
	}

	pm := config.NewProcessorManager()
	pm.StartHealthCheck()

	// Initialize Worker Pool
	workers.InitGlobalWorkerPool()
	defer workers.StopGlobalWorkerPool()

	routes := router.SetupRoutes(logger, pm)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: routes,
	}

	defer server.Close()

	logger.Infof("API listening at %s", server.Addr)
	server.ListenAndServe()
}
