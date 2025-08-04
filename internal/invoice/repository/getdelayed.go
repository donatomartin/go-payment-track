package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) GetDelayedInvoices(ctx context.Context) ([]invoice.Invoice, error) {
	query := `
		SELECT
		 invoices.id,
		 invoices.customer_name,
		 invoices.amount_due,
		 invoices.due_date,
		 COALESCE(SUM(payments.amount),0) as total_paid
		FROM invoices
		LEFT JOIN payments ON invoices.id = payments.invoice_id
		WHERE due_date < CURRENT_TIMESTAMP
		GROUP BY invoices.id, invoices.customer_name, invoices.amount_due
		HAVING total_paid < amount_due
	`
	rows, err := r.db.QueryContext(ctx, query)
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
