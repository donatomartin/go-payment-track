package repository

import (
	"app/internal/payment"
	"context"
)

func (r *PaymentRepository) UpdatePayment(ctx context.Context, payment *payment.Payment) error {
	_, err := r.db.ExecContext(ctx, "UPDATE payments SET invoice_id = ?, amount = ?, date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		payment.InvoiceID,
		payment.Amount,
		payment.Date,
		payment.ID,
	)

	return err
}
