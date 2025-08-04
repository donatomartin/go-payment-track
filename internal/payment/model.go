package payment

import (
	"time"
)

type Payment struct {
	ID         int
	InvoiceID  string
	ClientName string
	Amount     float64
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
