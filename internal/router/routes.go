package router

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/handlers"
)

func SetupRoutes(logger *config.Logger) http.Handler {
	server := http.NewServeMux()

	server.HandleFunc("POST /payments", handlers.PaymentHandler)
	server.HandleFunc("GET /payments-summary", handlers.PaymentSummaryHandler)
	server.HandleFunc("GET /payments/", handlers.PaymentDetailsHandler)

	logger.Info("All routes loaded!!!")

	return server
}
