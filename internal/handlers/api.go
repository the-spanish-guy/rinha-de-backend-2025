package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/helpers"
	"rinha-de-backend-2025/internal/logger"
	"rinha-de-backend-2025/internal/types"
	"rinha-de-backend-2025/internal/workers"
	"time"
)

var log = logger.GetLogger("[HANDLER]")

type Handler struct {
	processorManager *config.ProcessorManager
}

func HandleHandler(pm *config.ProcessorManager) *Handler {
	return &Handler{
		processorManager: pm,
	}
}

func (h *Handler) PaymentHandler(w http.ResponseWriter, r *http.Request) {
	readBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when trying read Body", http.StatusBadRequest)
	}

	// Parse request para Worker Pool
	var payment types.PaymentsRequest
	if err := json.Unmarshal(readBody, &payment); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Enviar para worker pool
	workers.GlobalWorkerPool.AddJob(payment)

	w.WriteHeader(http.StatusAccepted)
}

// este endpoint é só pra testes
func PaymentDetailsHandler(w http.ResponseWriter, r *http.Request) {
	HOST := os.Getenv("PROCESSOR_DEFAULT_URL")

	id := r.URL.Path[len("/payments/"):]

	res, errPayments := http.Get(HOST + "/payments/" + id)
	if errPayments != nil {
		fmt.Printf("Error on POST /payments %s", errPayments)
	}

	formattedResponse, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Error when trying read response body", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(formattedResponse)
}

func PaymentSummaryHandler(w http.ResponseWriter, r *http.Request) {
	var from, to *time.Time
	var query string
	var args []interface{}
	ctx := context.Background()

	if fromStr := r.URL.Query().Get("from"); fromStr != "" {
		if parsed, err := helpers.ParseFlexibleDateTime(fromStr); err == nil {
			from = &parsed
		}
	}

	if toStr := r.URL.Query().Get("to"); toStr != "" {
		if parsed, err := helpers.ParseFlexibleDateTime(toStr); err == nil {
			to = &parsed
		}
	}

	baseQuery := `
		select
			p.processor,
			count(1) as total_request,
			sum(p.amount) as total_amount
		from
			payments p
			WHERE p.status = 'SETTLED'
		`

	if from != nil && to != nil {
		query = baseQuery + `
			AND p.requested_at between $1 and $2
		`
		args = []interface{}{*from, *to}
	} else {
		query = baseQuery + `
			AND p.requested_at > NOW() - INTERVAL '1 year'
		`
	}

	query += `
		GROUP BY
			p.processor
		ORDER BY
			p.processor
	`

	rows, err := db.PGDB.Query(ctx, query, args...)

	if err != nil {
		log.Errorf("query execution failed: %v", err)
		http.Error(w, "query execution failed", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	response := types.PaymentsSummaryResponse{
		Default: types.SummaryResponse{
			TotalRequest: 0,
			TotalAmount:  0.0,
		},
		Fallback: types.SummaryResponse{
			TotalRequest: 0,
			TotalAmount:  0.0,
		},
	}
	for rows.Next() {
		var processor string
		var totalRequest int64
		var totalAmount float64

		if err := rows.Scan(&processor, &totalRequest, &totalAmount); err != nil {
			log.Errorf("failed to scan row: %v", err)
			http.Error(w, "failed to process results", http.StatusInternalServerError)
			return
		}

		summary := types.SummaryResponse{
			TotalRequest: int(totalRequest),
			TotalAmount:  totalAmount,
		}

		switch processor {
		case "DEFAULT":
			response.Default = summary
		case "FALLBACK":
			response.Fallback = summary
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ProcessorStatusHandler(w http.ResponseWriter, r *http.Request) {
	activeProcessor, _ := h.processorManager.GetActiveProcessor()

	response := map[string]interface{}{
		"activeProcessor": activeProcessor,
		"timestamp":       time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) PurgePayments(w http.ResponseWriter, r *http.Request) {
	db.GetDB().Exec(context.Background(), "DELETE FROM payments")
	w.WriteHeader(http.StatusOK)
}
