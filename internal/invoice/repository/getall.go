package repository

import (
	"app/internal/invoice"
	"context"
	"fmt"
)

func (r *InvoiceRepository) GetAll(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]invoice.Invoice, error) {
	validSortBy := map[string]bool{
		"id":            true,
		"customer_name": true,
		"amount_due":    true,
		"payment_mean":  true,
		"invoice_date":  true,
		"due_date":      true,
		"created_at":    true,
		"updated_at":    true,
		"total_paid":    true,
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

	query := fmt.Sprintf(`
		SELECT 
			invoices.*,
			COALESCE(SUM(payments.amount), 0) AS total_paid
		FROM invoices 
		LEFT JOIN payments ON invoices.id = payments.invoice_id
		GROUP BY invoices.id
		ORDER BY %s %s 
		LIMIT ? OFFSET ?`, sortBy, sortDir)
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
			&invoice.TotalPaid,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
