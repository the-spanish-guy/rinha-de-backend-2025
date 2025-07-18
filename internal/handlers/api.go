package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/messaging/nats"
	"rinha-de-backend-2025/internal/types"
	"time"

	"github.com/redis/go-redis/v9"
)

var logger = config.GetLogger("handler")

type Handler struct {
	publisher *nats.Publisher
}

func HandleHandler(p *nats.Publisher) *Handler {
	return &Handler{
		publisher: p,
	}
}

func (h *Handler) PaymentHandler(w http.ResponseWriter, r *http.Request) {
	HOST := healthcheck()

	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when trying read Body", http.StatusBadRequest)
	}

	res, errPayments := http.Post(HOST+"/payments", "application/json", bytes.NewReader(readBody))
	if errPayments != nil {
		logger.Errorf("Error on POST /payments %s", errPayments)
		// TODO: implementar retry
	}

	formattedResponse, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error when trying read response body", http.StatusBadRequest)
	}

	// requestedAt := time.Now().Unix()
	// TODO: remover o set do redis
	requestedAt, _ := time.Parse(time.RFC3339, "2025-07-16T14:20:00Z")
	err = db.DB.ZAdd(db.Ctx, "rinha-payments", redis.Z{
		Member: string(readBody),
		Score:  float64(requestedAt.Unix()),
	}).Err()
	if err != nil {
		logger.Errorf("Error on redis insert: %v", err)
	}

	var paymentRequest types.PaymentsRequest
	if err := json.Unmarshal(readBody, &paymentRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	}

	body := types.PaymentsRequest{
		CorrelationId: paymentRequest.CorrelationId,
		Amount:        paymentRequest.Amount,
	}

	msg, err := json.Marshal(body)
	if err != nil {
		// se a mensagem não foi publicada paciência.
		logger.Errorf("Invalid JSON format: %v", err)
	}

	h.publisher.PublishMessage(types.NewMessage(string(msg)))

	w.WriteHeader(http.StatusAccepted)
	w.Write(formattedResponse)
}

// este endpoint é só pra testes
func PaymentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	HOST := healthcheck()

	id := r.URL.Path[len("/payments/"):]

	res, errPayments := http.Get(HOST + "/payments/" + id)
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
	var from, to *time.Time

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if parsed, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = &parsed
		}
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if parsed, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = &parsed
		}
	}

	var min, max string

	if from != nil && to != nil {
		min = fmt.Sprintf("%d", from.Unix())
		max = fmt.Sprintf("%d", to.Unix())
	} else {
		min = "-inf"
		max = "+inf"
	}

	result, err := db.DB.ZRangeByScore(db.Ctx, "rinha-payments", &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		http.Error(w, "Error when trying read data", http.StatusBadRequest)
	}

	var totalRequests int = 0
	var totalAmount int
	for _, eventStr := range result {
		var event types.PaymentsRequest
		if err := json.Unmarshal([]byte(eventStr), &event); err != nil {
			continue
		}
		totalRequests++
		totalAmount += int(event.Amount)
	}
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"totalRequests": totalRequests,
		"totalAmount":   totalAmount,
	})
}

func healthcheck() string {
	// TODO: implementar a validacao para escolher entre o default ou fallback
	DEFAULT_HOST := os.Getenv("PROCESSOR_DEFAULT_URL")
	resp, err := http.Get(DEFAULT_HOST + "/payments/service-health")
	if err != nil {
		logger.Errorf("error service-health %v \n", err)
	}

	bodyHealth, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error %v \n", err)
	}
	logger.Infof("resp %s", bodyHealth)
	return DEFAULT_HOST
}
