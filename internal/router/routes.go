package router

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/handlers"
	"rinha-de-backend-2025/internal/handlers/middleware"
	"rinha-de-backend-2025/internal/logger"
	"rinha-de-backend-2025/internal/messaging/nats"
)

func SetupRoutes(logger *logger.Logger, pub *nats.Publisher) http.Handler {
	server := http.NewServeMux()

	handler := handlers.HandleHandler(pub)

	server.HandleFunc("POST /payments", handler.PaymentHandler)
	server.HandleFunc("GET /payments-summary", handlers.PaymentSummaryHandler)
	server.HandleFunc("GET /payments/", handlers.PaymentDetailsHandler)

	md := middleware.Logging(logger)(server)

	logger.Info("All routes loaded!!!")

	return md
}
