package main

import (
	"fmt"
	"net/http"
	"rinha-de-backend-2025/core/handlers"
)

func main() {
	fmt.Println("Teste")
	server := http.NewServeMux()
	server.HandleFunc("POST /payments", handlers.PaymentHandler)
	server.HandleFunc("POST /payments-summary", handlers.PaymentSummaryHandler)

	http.ListenAndServe("localhost:8080", server)
}
