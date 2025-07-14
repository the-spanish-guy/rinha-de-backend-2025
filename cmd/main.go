package main

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/router"
)

func main() {
	logger := config.GetLogger("main")

	logger.Info("Starting project")

	db.StartDB()
	handler := router.SetupRoutes(logger)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: handler,
	}

	logger.Infof("API listening at %s", server.Addr)
	server.ListenAndServe()
}
