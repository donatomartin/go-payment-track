package repository

import (
	"app/internal/payment"
	"context"
	"fmt"
)

func (r *PaymentRepository) GetAll(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]payment.Payment, error) {

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

	var payments []payment.Payment
	for rows.Next() {
		var payment payment.Payment
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
