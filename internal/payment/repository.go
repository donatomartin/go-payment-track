package payment

import (
	"context"
	"database/sql"
	"fmt"
)

type PaymentRepository interface {
	GetAll(ctx context.Context) ([]Payment, error)
	GetPaged(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Payment, error)
	AddPayment(ctx context.Context, payment *Payment) error
	CreateBatch(ctx context.Context, payments []Payment) error
	UpdatePayment(ctx context.Context, payment *Payment) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]Payment, error) {
	rows, err := (*r).db.QueryContext(ctx, "SELECT * FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
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

	query := fmt.Sprintf(`
    SELECT 
        payments.*, 
        invoices.customer_name 
    FROM payments 
    JOIN invoices ON payments.invoice_id = invoices.ID 
    ORDER BY %s %s 
    LIMIT ? OFFSET ?`, sortBy, sortDir)

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
			&payment.ClientName,
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

func (r *paymentRepository) CreateBatch(ctx context.Context, payments []Payment) error {
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

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *Payment) error {
	_, err := r.db.ExecContext(ctx, "UPDATE payments SET invoice_id = ?, amount = ?, date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		payment.InvoiceID,
		payment.Amount,
		payment.Date,
		payment.ID,
	)

	return err
}
