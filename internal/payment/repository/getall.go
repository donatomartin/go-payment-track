package repository

import (
	"app/internal/payment"
	"context"
)

func (r *PaymentRepository) GetAll(ctx context.Context) ([]payment.Payment, error) {
	rows, err := (*r).db.QueryContext(ctx, "SELECT * FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []payment.Payment
	for rows.Next() {
		var payment payment.Payment
		if err := rows.Scan(
			&payment.InvoiceID,
			&payment.ID,
			&payment.Amount,
			&payment.Date,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil

}
