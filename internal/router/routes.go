package router

import (
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/handlers"
	"rinha-de-backend-2025/internal/handlers/middleware"
	"rinha-de-backend-2025/internal/logger"
)

func SetupRoutes(logger *logger.Logger, processorManager *config.ProcessorManager) http.Handler {
	server := http.NewServeMux()

	handler := handlers.HandleHandler(processorManager)

	server.HandleFunc("POST /payments", handler.PaymentHandler)
	server.HandleFunc("GET /payments-summary", handlers.PaymentSummaryHandler)
	server.HandleFunc("POST /admin/purge-payments", handler.PurgePayments)

	// endpoints usado para testes, não necessarios para a rinha
	server.HandleFunc("GET /processors/status", handler.ProcessorStatusHandler)
	server.HandleFunc("GET /payments/", handlers.PaymentDetailsHandler)

	md := middleware.Logging(logger)(server)

	logger.Info("All routes loaded!!!")

	return md
}
