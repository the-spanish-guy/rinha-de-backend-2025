package main

import (
	"net/http"
	"rinha-de-backend-2025/core/handlers"
)

func main() {
	server := http.NewServeMux()
	server.HandleFunc("POST /payments", handlers.PaymentHandler)
	server.HandleFunc("GET /payments-summary", handlers.PaymentSummaryHandler)

	http.ListenAndServe("0.0.0.0:8080", server)
}
