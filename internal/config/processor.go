package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/logger"
	"rinha-de-backend-2025/internal/types"
	"time"
)

var (
	log             = logger.GetLogger("PROCESSOR")
	processorConfig = &ProcessorConfig{
		DefaultURL:          os.Getenv("PROCESSOR_DEFAULT_URL"),
		FallbackURL:         os.Getenv("PROCESSOR_FALLBACK_URL"),
		Timeout:             5,
		HealthCheckInterval: 8,
		CacheTTL:            8,
	}
)

type ProcessorConfig struct {
	DefaultURL          string
	FallbackURL         string
	Timeout             int // healthcheck timeout in seconds
	HealthCheckInterval int // healthcheck interval in seconds
	CacheTTL            int // cache TTL in seconds
}

type ProcessorManager struct {
	config    *ProcessorConfig
	ctx       context.Context
	activeURL string
	lastCheck time.Time
}

type healthResult struct {
	service string
	data    []byte
	err     error
}

func NewProcessorManager() *ProcessorManager {
	return &ProcessorManager{
		config:    processorConfig,
		ctx:       context.Background(),
		activeURL: processorConfig.DefaultURL,
		lastCheck: time.Now(),
	}
}

func (pm *ProcessorManager) GetActiveProcessor() (string, error) {
	var defaultPaylaod, fallbackPayload types.HealtchCheckResponse
	if err := pm.getCachedProcessor("health:default", &defaultPaylaod); err != nil {
		log.Errorf("an error occurred while decode host from redis")
		return "", nil
	}
	if err := pm.getCachedProcessor("health:fallback", &fallbackPayload); err != nil {
		log.Errorf("an error occurred while decode host from redis")
		return "", nil
	}

	var HOST = os.Getenv("PROCESSOR_DEFAULT_URL")
	if defaultPaylaod.Failing || defaultPaylaod.MinResponseTime > fallbackPayload.MinResponseTime {
		HOST = os.Getenv("PROCESSOR_FALLBACK_URL")
	}

	return HOST, nil
}

func (pm *ProcessorManager) getCachedProcessor(key string, result interface{}) error {
	if db.DB == nil {
		return fmt.Errorf("")
	}

	jsonData, err := db.DB.Get(pm.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonData), result)
}

func (pm *ProcessorManager) updatProcessorCache(key string, data []byte) {
	if db.DB == nil {
		return
	}

	err := db.DB.Set(pm.ctx, key, data, time.Duration(pm.config.CacheTTL)*time.Second).Err()
	if err != nil {
		log.Errorf("Error saving processor to cache: %v", err)
	}
}

func (pm *ProcessorManager) healthCheck() {
	client := &http.Client{
		Timeout: time.Duration(pm.config.Timeout) * time.Second,
	}
	defaultURL := os.Getenv("PROCESSOR_DEFAULT_URL")
	fallbackURL := os.Getenv("PROCESSOR_FALLBACK_URL")

	results := make(chan healthResult, 2)

	go pm.checkServiceHealth(client, defaultURL, "default", results)
	go pm.checkServiceHealth(client, fallbackURL, "fallback", results)

	for i := 0; i < 2; i++ {
		result := <-results
		if result.err != nil {
			log.Errorf("Health check failed for %s: %v", result.service, result.err)
			pm.updatProcessorCache(fmt.Sprintf("health:%s:error", result.service), []byte(`{"status":"error"}`))
		} else {
			pm.updatProcessorCache(fmt.Sprintf("health:%s", result.service), result.data)
		}
	}
}

func (pm *ProcessorManager) checkServiceHealth(client *http.Client, baseURL, serviceName string, results chan<- healthResult) {
	response, err := client.Get(baseURL + "/payments/service-health")
	if err != nil {
		results <- healthResult{service: serviceName, err: err}
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		results <- healthResult{
			service: serviceName,
			err:     fmt.Errorf("HTTP %d", response.StatusCode),
		}
		return
	}

	var payload types.HealtchCheckResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		results <- healthResult{service: serviceName, err: err}
		return
	}

	data, err := json.Marshal(&payload)
	if err != nil {
		results <- healthResult{service: serviceName, err: err}
		return
	}

	results <- healthResult{service: serviceName, data: data}
}

func (pm *ProcessorManager) StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(time.Duration(pm.config.HealthCheckInterval) * time.Second)
		defer ticker.Stop()
		log.Infof("Starting automatic health check loop (every %ds)", pm.config.HealthCheckInterval)

		for {
			select {
			case <-ticker.C:
				pm.healthCheck()
				log.Infof("Active processor: %s", pm.activeURL)
			case <-pm.ctx.Done():
				log.Info("Health check loop stopped")
				return
			}
		}
	}()
}
