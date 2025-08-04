package invoice

import (
	"time"
)

type Invoice struct {
	ID           string
	CustomerName string
	AmountDue    float64
	PaymentMean  string
	InvoiceDate  time.Time
	DueDate      time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	TotalPaid    float64
}
