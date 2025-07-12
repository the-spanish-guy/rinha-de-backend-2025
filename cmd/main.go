package main

import (
	"fmt"
	"net/http"
	"rinha-de-backend-2025/endpoints"
)

func main() {
	fmt.Println("Teste")
	server := http.NewServeMux()
	server.HandleFunc("POST /payments", endpoints.PaymentHandler)
	server.HandleFunc("POST /payments-summary", endpoints.PaymentSummaryHandler)

	http.ListenAndServe("localhost:8080", server)
}
