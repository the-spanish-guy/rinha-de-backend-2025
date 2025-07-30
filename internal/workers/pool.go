package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"rinha-de-backend-2025/internal/config"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/logger"
	"rinha-de-backend-2025/internal/types"
	"strings"
	"sync"
	"time"
)

var log = logger.GetLogger("[WORKERS]")

type WorkerPool struct {
	jobQueue    chan types.PaymentsRequest
	workers     []*Worker
	wg          sync.WaitGroup
	workerCount int
	isRunning   bool
	mu          sync.RWMutex
}

type Worker struct {
	id       int
	jobQueue chan types.PaymentsRequest
	quit     chan bool
	pm       *config.ProcessorManager
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan types.PaymentsRequest, 1000),
		workerCount: workerCount,
		workers:     make([]*Worker, workerCount),
	}
}

func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.isRunning {
		return
	}

	log.Infof("Starting worker pool with %d workers", wp.workerCount)

	pm := config.NewProcessorManager()

	for i := 0; i < wp.workerCount; i++ {
		worker := &Worker{
			id:       i,
			jobQueue: wp.jobQueue,
			quit:     make(chan bool),
			pm:       pm,
		}

		wp.workers[i] = worker
		wp.wg.Add(1)
		go worker.start(&wp.wg)
	}

	wp.isRunning = true
	log.Info("Worker pool started successfully")
}

func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.isRunning {
		return
	}

	log.Info("Stopping worker pool...")

	close(wp.jobQueue)

	for _, worker := range wp.workers {
		close(worker.quit)
	}

	wp.wg.Wait()
	wp.isRunning = false

	log.Info("Worker pool stopped")
}

func (wp *WorkerPool) AddJob(payment types.PaymentsRequest) bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	if !wp.isRunning {
		return false
	}

	select {
	case wp.jobQueue <- payment:
		return true
	default:
		log.Errorf("Job queue is full, dropping payment")
		return false
	}
}

func (w *Worker) start(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Infof("Worker %d started", w.id)

	for {
		select {
		case payment, ok := <-w.jobQueue:
			if !ok {
				log.Infof("Worker %d stopping - job queue closed", w.id)
				return
			}
			w.processPayment(payment)

		case <-w.quit:
			log.Infof("Worker %d stopping - quit signal received", w.id)
			return
		}
	}
}

func (w *Worker) processPayment(payment types.PaymentsRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	payment.RequestedAt = time.Now()

	activeHost, err := w.pm.GetActiveProcessor()
	if err != nil {
		log.Errorf("Worker %d: Error getting active processor: %v", w.id, err)
		return
	}

	processorType := getProcessorType(activeHost)

	paymentDB := types.Payments{
		CorrelationId: payment.CorrelationId,
		Amount:        payment.Amount,
		RequestedAt:   payment.RequestedAt,
		Status:        "PENDING",
		Processor:     processorType,
	}

	pgdb := db.GetDB()
	if pgdb == nil {
		log.Errorf("Worker %d: Database connection is nil", w.id)
		return
	}

	_, err = pgdb.Exec(ctx,
		`INSERT INTO payments (correlation_id, amount, status, processor, requested_at) 
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (correlation_id) DO UPDATE SET
		 	updated_at = NOW()
`,
		paymentDB.CorrelationId, paymentDB.Amount, paymentDB.Status, paymentDB.Processor, paymentDB.RequestedAt)

	if err != nil {
		log.Errorf("Worker %d: Failed to insert payment: %v", w.id, err)
		return
	}

	maxRetries := 2
	var finalStatus string = "FAILED"

	client := &http.Client{
		Timeout: 2000 * time.Millisecond,
		Transport: &http.Transport{
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 25,
			IdleConnTimeout:     10 * time.Second,
			DisableCompression:  true,
			DisableKeepAlives:   false,
		},
	}

	payloadForProcessor := struct {
		CorrelationId string  `json:"correlationId"`
		Amount        float64 `json:"amount"`
	}{
		CorrelationId: payment.CorrelationId,
		Amount:        payment.Amount,
	}
	requestBody, _ := json.Marshal(payloadForProcessor)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(100*attempt*attempt) * time.Millisecond
			time.Sleep(delay)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", activeHost+"/payments", bytes.NewReader(requestBody))
		if err != nil {
			log.Errorf("Worker %d: Failed to create request (attempt %d): %v", w.id, attempt+1, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)

		if err != nil {
			log.Errorf("Worker %d: HTTP request failed (attempt %d): %v", w.id, attempt+1, err)
			continue
		}

		defer res.Body.Close()

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			finalStatus = "SETTLED"
			break
		}

		if res.StatusCode == 429 || res.StatusCode == 503 {
			continue
		}

		break
	}

	result, err := pgdb.Exec(ctx,
		"UPDATE payments SET status = $1, updated_at = NOW() WHERE correlation_id = $2 AND status = 'PENDING'",
		finalStatus, paymentDB.CorrelationId)

	if err != nil {
		log.Errorf("Worker %d: Failed to update payment status: %v", w.id, err)
		return
	}

	if result.RowsAffected() == 0 {
		log.Errorf("Worker %d: Payment %s not updated - possibly already processed", w.id, paymentDB.CorrelationId)
	}
}

func getProcessorType(host string) string {
	if strings.Contains(host, "fallback") {
		return "FALLBACK"
	}
	return "DEFAULT"
}

var GlobalWorkerPool *WorkerPool

func InitGlobalWorkerPool() {
	GlobalWorkerPool = NewWorkerPool(8)
	GlobalWorkerPool.Start()
}

func StopGlobalWorkerPool() {
	if GlobalWorkerPool != nil {
		GlobalWorkerPool.Stop()
	}
}
