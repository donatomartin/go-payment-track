package repository

import "context"

func (r *PaymentRepository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM payments")
	return err
}
