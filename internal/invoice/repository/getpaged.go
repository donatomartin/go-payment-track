package repository

import (
	"app/internal/invoice"
	"context"
	"fmt"
)

func (r *InvoiceRepository) GetPaged(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]invoice.Invoice, error) {
	validSortBy := map[string]bool{
		"id":            true,
		"customer_name": true,
		"amount_due":    true,
		"payment_mean":  true,
		"invoice_date":  true,
		"due_date":      true,
		"created_at":    true,
		"updated_at":    true,
	}

	if !validSortBy[sortBy] {
		return nil, fmt.Errorf("invalid sort by field: %s", sortBy)
	}

	validSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortDir[sortDir] {
		return nil, fmt.Errorf("invalid sort direction: %s", sortDir)
	}

	query := "SELECT * FROM invoices ORDER BY " + sortBy + " " + sortDir + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []invoice.Invoice
	for rows.Next() {
		var invoice invoice.Invoice
		if err := rows.Scan(
			&invoice.ID,
			&invoice.CustomerName,
			&invoice.AmountDue,
			&invoice.PaymentMean,
			&invoice.InvoiceDate,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
