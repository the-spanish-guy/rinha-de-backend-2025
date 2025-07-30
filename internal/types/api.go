package types

import "time"

type HealtchCheckResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

type PaymentsRequest struct {
	CorrelationId string    `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
}

type SummaryResponse struct {
	TotalRequest int     `json:"totalRequests"`
	TotalAmount  float64 `json:"totalAmount"`
}

type PaymentsSummaryResponse struct {
	Default  SummaryResponse `json:"default"`
	Fallback SummaryResponse `json:"fallback"`
}

type ProcessorHealth struct {
	URL          string    `json:"url"`
	IsHealthy    bool      `json:"isHealthy"`
	ResponseTime int64     `json:"responseTime"`
	LastCheck    time.Time `json:"lastCheck"`
	ErrorCount   int       `json:"errorCount"`
	SuccessCount int       `json:"successCount"`
}

type ProcessorConfig struct {
	DefaultURL  string `json:"defaultUrl"`
	FallbackURL string `json:"fallbackUrl"`
	MaxErrors   int    `json:"maxErrors"`
	Timeout     int    `json:"timeout"`
}
