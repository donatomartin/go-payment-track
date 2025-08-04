package repository

import (
	"context"
)

func (r *InvoiceRepository) GetDelayedInvoicesAmount(ctx context.Context) (float64, error) {
	query := `
		SELECT SUM(amount_due - COALESCE(total_paid, 0)) FROM (
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
		) AS delayed_invoices
	`

	var total float64
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
