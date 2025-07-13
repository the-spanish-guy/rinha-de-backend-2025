package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rinha-de-backend-2025/core/db"
	"rinha-de-backend-2025/core/types"
	"time"

	"github.com/redis/go-redis/v9"
)

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	HOST := healthcheck()

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when trying read Body", http.StatusBadRequest)
	}

	var paymentRequest types.Payments
	if err := json.Unmarshal(readBody, &paymentRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	}

	res, errPayments := http.Post(HOST+"/payments", "application/json", bytes.NewReader(readBody))
	if errPayments != nil {
		fmt.Printf("Error on POST /payments %s", errPayments)
	}

	formattedResponse, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error when trying read response body", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(formattedResponse)
}

func PaymentSummaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deu bom")
	w.WriteHeader(http.StatusAccepted)
	// w.Write(resp)
}

func healthcheck() string {
	// TODO: implementar a validacao para escolher entre o default ou fallback
	DEFAULT_HOST := os.Getenv("PROCESSOR_DEFAULT_URL")
	resp, err := http.Get(DEFAULT_HOST + "/payments/service-health")
	if err != nil {
		fmt.Printf("error %v \n", err)
	}

	bodyHealth, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error %v \n", err)
	}
	fmt.Printf("resp %s", bodyHealth)
	return DEFAULT_HOST
}
