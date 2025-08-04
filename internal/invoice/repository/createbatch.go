package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) CreateBatch(ctx context.Context, invoices []invoice.Invoice) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO invoices (id, customer_name, amount_due, payment_mean, invoice_date, due_date) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, invoice := range invoices {
		if _, err := stmt.Exec(
			invoice.ID,
			invoice.CustomerName,
			invoice.AmountDue,
			invoice.PaymentMean,
			invoice.InvoiceDate,
			invoice.DueDate,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
