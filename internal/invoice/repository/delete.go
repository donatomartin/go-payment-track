package repository

import (
	"context"
)

func (r *InvoiceRepository) DeleteInvoice(ctx context.Context, id int) error {
	query := "DELETE FROM invoices WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}
