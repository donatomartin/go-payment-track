package repository

import "context"

func (r *InvoiceRepository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM invoices")
	return err
}
