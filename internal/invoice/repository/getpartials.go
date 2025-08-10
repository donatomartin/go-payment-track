package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) GetPartialInvoices(ctx context.Context, offset, limit int) ([]invoice.Invoice, error) {
	query := `
                SELECT
                        invoices.*,
                        COALESCE(SUM(payments.amount),0) as total_paid
                FROM invoices
                JOIN payments ON invoices.id = payments.invoice_id
                GROUP BY invoices.id
                HAVING total_paid != amount_due
                ORDER BY invoices.created_at DESC
                LIMIT ? OFFSET ?
        `
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
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
