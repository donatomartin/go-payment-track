package payment

type Payment struct {
	ID        int
	InvoiceID string
	Amount    float64
	Date      string
	CreatedAt string
	UpdatedAt string
}
