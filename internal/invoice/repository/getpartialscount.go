package repository

import (
	"context"
)

func (r *InvoiceRepository) GetPartialInvoicesCount(ctx context.Context) (int, error) {
	query := `
	SELECT COUNT(*) FROM (
		SELECT
		 invoices.id,
		 invoices.customer_name,
		 invoices.amount_due,
		 SUM(payments.amount) as total_paid
		FROM invoices
		JOIN payments ON invoices.id = payments.invoice_id
		GROUP BY invoices.id, invoices.customer_name, invoices.amount_due
		HAVING total_paid != amount_due
	)
	`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
