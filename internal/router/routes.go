package router

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/handlers"
	"rinha-de-backend-2025/internal/handlers/middleware"
)

func SetupRoutes(logger *config.Logger) http.Handler {
	server := http.NewServeMux()

	server.HandleFunc("POST /payments", handlers.PaymentHandler)
	server.HandleFunc("GET /payments-summary", handlers.PaymentSummaryHandler)
	server.HandleFunc("GET /payments/", handlers.PaymentDetailsHandler)

	handler := middleware.Logging(logger)(server)

	logger.Info("All routes loaded!!!")

	return handler
}
