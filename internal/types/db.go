package types

import "time"

type Payments struct {
	Id            string    `json:"id"`
	CorrelationId string    `json:"correlationId" db:"correlation_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Status        string    `json:"status" db:"status"`
	Processor     string    `json:"processor" db:"processor"`
	RequestedAt   time.Time `json:"requestedAt" db:"requested_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}
