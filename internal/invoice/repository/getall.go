package repository

import (
	"app/internal/invoice"
	"context"
)

func (r *InvoiceRepository) GetAll(ctx context.Context) ([]invoice.Invoice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM invoices")
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
