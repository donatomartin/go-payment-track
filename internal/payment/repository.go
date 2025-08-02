package payment

import (
	"context"
	"database/sql"
	"fmt"
)

type PaymentRepository interface {
	GetAll(ctx context.Context) ([]Payment, error)
	GetPaged(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]Payment, error) {
	rows, err := (*r).db.QueryContext(ctx, "SELECT id, amount FROM payments")
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

	validSortBy := map[string]bool{
		"id":         true,
		"invoice_id": true,
		"amount":     true,
		"date":       true,
		"created_at": true,
		"updated_at": true,
	}

	validSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortBy[sortBy] {
		return nil, fmt.Errorf("invalid sortBy field: %s", sortBy)
	}

	if !validSortDir[sortDir] {
		return nil, fmt.Errorf("invalid sortDir value: %s", sortDir)
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset cannot be negative: %d", offset)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit must be at least 1: %d", limit)
	}

	query := fmt.Sprintf("SELECT * FROM payments ORDER BY %s %s LIMIT ? OFFSET ?", sortBy, sortDir)

	rows, err := (*r).db.QueryContext(ctx, query, limit, offset)
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

func (r *paymentRepository) AddPayment(ctx context.Context, payment *Payment) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO payments (invoice_id, amount, date) VALUES (?, ?, ?)",
		payment.InvoiceID,
		payment.Amount,
		payment.Date,
	)

	return err

}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *Payment) error {
	_, err := r.db.ExecContext(ctx, "UPDATE payments SET invoice_id = ?, amount = ?, date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		payment.InvoiceID,
		payment.Amount,
		payment.Date,
		payment.ID,
	)

	return err
}
