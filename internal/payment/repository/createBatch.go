package repository

import (
	"app/internal/payment"
	"context"
)

func (r *PaymentRepository) CreateBatch(ctx context.Context, payments []payment.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	query := "INSERT INTO payments (invoice_id, amount, date) VALUES "
	values := make([]any, 0, len(payments)*3)

	for i, payment := range payments {
		if i > 0 {
			query += ", "
		}
		query += "(?, ?, ?)"
		values = append(values, payment.InvoiceID, payment.Amount, payment.Date)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	return err
}
