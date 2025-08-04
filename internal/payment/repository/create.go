package repository

import (
	"app/internal/payment"
	"context"
)

func (r *PaymentRepository) AddPayment(ctx context.Context, payment *payment.Payment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO payments (invoice_id, amount, date) VALUES (?, ?, ?)",
		payment.InvoiceID,
		payment.Amount,
		payment.Date,
	)

	return err

}
