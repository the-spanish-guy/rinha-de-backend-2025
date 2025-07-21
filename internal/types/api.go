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
