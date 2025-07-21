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
	TotalRequest string `json:"totalRequests"`
	TotalAmount  string `json:"totalAmount"`
}

type PaymentsSummaryResponse struct {
	Default  SummaryResponse `json:"default"`
	Fallback SummaryResponse `json:"fallback"`
}
