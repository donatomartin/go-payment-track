package invoice

import (
	"context"
	"database/sql"
	"fmt"
)

type InvoiceRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

func (r *InvoiceRepository) GetAll(ctx context.Context) ([]Invoice, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM invoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
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

func (r *InvoiceRepository) GetPaged(ctx context.Context, sortBy, sortDir string, offset, limit int) ([]Invoice, error) {
	validSortBy := map[string]bool{
		"id":            true,
		"customer_name": true,
		"amount_due":    true,
		"payment_mean":  true,
		"invoice_date":  true,
		"due_date":      true,
		"created_at":    true,
		"updated_at":    true,
	}

	if !validSortBy[sortBy] {
		return nil, fmt.Errorf("invalid sort by field: %s", sortBy)
	}

	validSortDir := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	if !validSortDir[sortDir] {
		return nil, fmt.Errorf("invalid sort direction: %s", sortDir)
	}

	query := "SELECT * FROM invoices ORDER BY " + sortBy + " " + sortDir + " LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
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

func (r *InvoiceRepository) AddInvoice(ctx context.Context, invoice *Invoice) error {
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

func (r *InvoiceRepository) CreateBatch(ctx context.Context, invoices []Invoice) error {
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

func (r *InvoiceRepository) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	query := "UPDATE invoices SET customer_name = ?, amount_due = ?, payment_mean = ?, invoice_date = ?, due_date = ? updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.Exec(query,
		invoice.CustomerName,
		invoice.AmountDue,
		invoice.PaymentMean,
		invoice.InvoiceDate,
		invoice.DueDate,
		invoice.ID,
	)
	return err
}

func (r *InvoiceRepository) DeleteInvoice(ctx context.Context, id int) error {
	query := "DELETE FROM invoices WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

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

func (r *InvoiceRepository) GetDelayedInvoices(ctx context.Context) ([]Invoice, error) {
	query := `
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
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(
			&invoice.ID,
			&invoice.CustomerName,
			&invoice.AmountDue,
			&invoice.PaymentMean,
			&invoice.InvoiceDate,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&invoice.TotalPaid,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (r *InvoiceRepository) GetPendingInvoicesCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) FROM (
			SELECT
			 invoices.id,
			 invoices.customer_name,
			 invoices.amount_due,
			 COALESCE(SUM(payments.amount),0) as total_paid
			FROM invoices
			LEFT JOIN payments ON invoices.id = payments.invoice_id
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

func (r *InvoiceRepository) GetPendingInvoices(ctx context.Context) ([]Invoice, error) {
	query := `
		SELECT
		 invoices.id,
		 invoices.customer_name,
		 invoices.amount_due,
		 COALESCE(SUM(payments.amount),0) as total_paid
		FROM invoices
		LEFT JOIN payments ON invoices.id = payments.invoice_id
		GROUP BY invoices.id, invoices.customer_name, invoices.amount_due
		HAVING total_paid < amount_due
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invoices []Invoice
	for rows.Next() {
		var invoice Invoice
		if err := rows.Scan(
			&invoice.ID,
			&invoice.CustomerName,
			&invoice.AmountDue,
			&invoice.PaymentMean,
			&invoice.InvoiceDate,
			&invoice.DueDate,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
			&invoice.TotalPaid,
		); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}
