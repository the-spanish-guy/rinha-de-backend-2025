package main

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/messaging"
	"rinha-de-backend-2025/internal/router"
)

func main() {
	logger := config.GetLogger("main")

	logger.Info("Starting project")
	db.StartDB()

	pub, sub := messaging.SetupMessaging()
	// Esse sub.Subscribe() talvez pudesse um método dentro da pasta de handlers
	sub.Subscribe()

	routes := router.SetupRoutes(logger, pub)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: routes,
	}

	logger.Infof("API listening at %s", server.Addr)
	server.ListenAndServe()
}
