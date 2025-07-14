package types

type HealtchCheckResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}

type Payments struct {
	CorrelationId string  `json:"correlationId"`
	Amount        float64 `json:"amount"`
}
