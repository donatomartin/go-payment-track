package database

import (
	"database/sql"
)

func GetSchema() string {
	return `
		DROP TABLE IF EXISTS payments;
		CREATE TABLE IF NOT EXISTS payments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				invoice_id TEXT NOT NULL,
				amount NUMERIC NOT NULL,
				date TIMESTAMP NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		DROP TABLE IF EXISTS invoices;
		CREATE TABLE IF NOT EXISTS invoices (
				id TEXT PRIMARY KEY,				
				customer_id INT NOT NULL,
				amount_due NUMERIC NOT NULL,
				payment_mean TEXT NOT NULL, 
				invoice_date TIMESTAMP NOT NULL,
				due_date TIMESTAMP,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
}

func ApplySchema(db *sql.DB) error {

	_, err := db.Exec(GetSchema())
	if err != nil {
		return err
	}

	return nil

}
