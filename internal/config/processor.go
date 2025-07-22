package config

import (
	"context"
	"net/http"
	"os"
	"rinha-de-backend-2025/internal/db"
	"rinha-de-backend-2025/internal/logger"
	"time"
)

var (
	log             = logger.GetLogger("PROCESSOR")
	processorConfig = &ProcessorConfig{
		DefaultURL:          os.Getenv("PROCESSOR_DEFAULT_URL"),
		FallbackURL:         os.Getenv("PROCESSOR_FALLBACK_URL"),
		Timeout:             2,
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

func NewProcessorManager() *ProcessorManager {
	return &ProcessorManager{
		config:    processorConfig,
		ctx:       context.Background(),
		activeURL: processorConfig.DefaultURL,
		lastCheck: time.Now(),
	}
}

func (pm *ProcessorManager) GetActiveProcessor() string {
	if cachedURL := pm.getCachedProcessor(); cachedURL != "" {
		log.Debugf("Cached processor: %s", cachedURL)
		pm.activeURL = cachedURL
		return cachedURL
	}

	pm.updateActiveProcessor()
	return pm.activeURL
}

func (pm *ProcessorManager) getCachedProcessor() string {
	if db.DB == nil {
		return ""
	}

	result, err := db.DB.Get(pm.ctx, "active_processor").Result()
	if err != nil {
		return ""
	}

	return result
}

func (pm *ProcessorManager) cacheProcessor(url string) {
	if db.DB == nil {
		return
	}

	err := db.DB.Set(pm.ctx, "active_processor", url, time.Duration(pm.config.CacheTTL)*time.Second).Err()
	if err != nil {
		log.Errorf("Error saving processor to cache: %v", err)
	}
}

func (pm *ProcessorManager) updateActiveProcessor() {
	if pm.isProcessorHealthy(pm.config.DefaultURL) {
		if pm.activeURL != pm.config.DefaultURL {
			log.Info("Changing to default processor")
		}
		pm.activeURL = pm.config.DefaultURL
		pm.lastCheck = time.Now()
		pm.cacheProcessor(pm.activeURL)
		return
	}

	if pm.isProcessorHealthy(pm.config.FallbackURL) {
		if pm.activeURL != pm.config.FallbackURL {
			log.Warning("Changing to fallback processor")
		}
		pm.activeURL = pm.config.FallbackURL
		pm.lastCheck = time.Now()
		pm.cacheProcessor(pm.activeURL)
		return
	}

	if pm.activeURL == "" {
		pm.activeURL = pm.config.DefaultURL
		log.Warning("No healthy processor, using default as last resort")
	}
	pm.lastCheck = time.Now()
	pm.cacheProcessor(pm.activeURL)
	log.Infof("Health check END - Selected: %s", pm.activeURL)
}

func (pm *ProcessorManager) isProcessorHealthy(url string) bool {
	client := &http.Client{
		Timeout: time.Duration(pm.config.Timeout) * time.Second,
	}

	// validar os campos failing e minResponseTime?
	resp, err := client.Get(url + "/payments/service-health")
	if err != nil {
		log.Errorf("Error checking health of processor %s: %v", url, err)
		return false
	}
	defer resp.Body.Close()

	isHealthy := resp.StatusCode == 200

	if !isHealthy {
		log.Warningf("Processor %s is not healthy: status=%d", url, resp.StatusCode)
	}

	return isHealthy
}

func (pm *ProcessorManager) StartHealthCheck() {
	go func() {
		ticker := time.NewTicker(time.Duration(pm.config.HealthCheckInterval) * time.Second)
		defer ticker.Stop()
		log.Infof("Starting automatic health check loop (every %ds)", pm.config.HealthCheckInterval)

		for {
			select {
			case <-ticker.C:
				pm.updateActiveProcessor()
				log.Infof("Active processor: %s", pm.activeURL)
			case <-pm.ctx.Done():
				log.Info("Health check loop stopped")
				return
			}
		}
	}()
}
