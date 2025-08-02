package repository

import (
	"context"
	"database/sql"
)

type Payment struct {
	ID        int
	InvoiceID string
	Amount    float64
	Date      string
	CreatedAt string
	UpdatedAt string
}

type PaymentRepository interface {
	GetAll(ctx context.Context) ([]Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]Payment, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, amount FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		if err := rows.Scan(&payment.ID, &payment.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil

}

func (r *paymentRepository) GetPaged(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Payment, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM payments ORDER BY $1 $2 OFFSET $3 LIMIT $4", sortBy, sortDir, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		if err := rows.Scan(
			&payment.ID,
			&payment.InvoiceID,
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
