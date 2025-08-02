package repository

import (
	"database/sql"
)

func ApplySchema(db *sql.DB) error {

	schema := `
    CREATE TABLE IF NOT EXISTS payments (
        id SERIAL PRIMARY KEY,
				invoice_id VARCHAR(255) NOT NULL,
        amount NUMERIC(10, 2) NOT NULL,
				date TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT now()
				updated_at TIMESTAMP DEFAULT now()
    );
		CREATE TABLE IF NOT EXISTS invoices (
				id VARCHAR(255) PRIMARY KEY,				
				customer_id INT NOT NULL,
				amount_due NUMERIC(10, 2) NOT NULL,
				payment_mean VARCHAR(50),
				invoice_date TIMESTAMP NOT NULL,
				due_date TIMESTAMP,
				created_at TIMESTAMP DEFAULT now()
				updated_at TIMESTAMP DEFAULT now()
		);
    `
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return nil

}
