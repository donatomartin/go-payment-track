package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) UpdateInvoice(ctx context.Context, invoice *invoice.Invoice) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE invoices SET customer_name = ?, amount_due = ?, payment_mean = ?, invoice_date = ?, due_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		invoice.CustomerName,
		invoice.AmountDue,
		invoice.PaymentMean,
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.ID,
	)
	return err
}
