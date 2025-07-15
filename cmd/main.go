package main

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/messaging/nats"
	"rinha-de-backend-2025/internal/router"
)

func main() {
	logger := config.GetLogger("main")

	logger.Info("Starting project")
	db.StartDB()

	// start nats server
	ns := nats.NewServer(logger)
	if err := ns.Start(); err != nil {
		logger.Fatal("An error occurred when NATS server starting", err.Error())
	}

	// create subscribe
	sub := nats.NewSubscriber()
	// connect sub to nats server
	if err := sub.Connect(); err != nil {
		logger.Fatal("An error occurred when create subscriber server", err.Error())
	}

	//create publisher
	pub := nats.NewPublisher()
	// connect pub to nats server
	if err := pub.Connect(); err != nil {
		logger.Fatal("An error occurred when create publisher server", err.Error())
	}

	handler := router.SetupRoutes(logger, pub)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: handler,
	}

	logger.Infof("API listening at %s", server.Addr)
	server.ListenAndServe()
}
