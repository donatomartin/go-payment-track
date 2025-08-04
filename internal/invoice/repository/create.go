package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) AddInvoice(ctx context.Context, invoice *invoice.Invoice) error {
	query := "INSERT INTO invoices (id, customer_name, amount_due, payment_mean, invoice_date, due_date) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.db.Exec(query,
		invoice.ID,
		invoice.CustomerName,
		invoice.AmountDue,
		invoice.PaymentMean,
		invoice.InvoiceDate,
		invoice.DueDate,
	)
	return err
}
