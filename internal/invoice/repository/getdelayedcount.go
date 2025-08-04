package repository

import (
	"context"
)

func (r *InvoiceRepository) GetDelayedInvoicesCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
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
		)
	`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil

}
