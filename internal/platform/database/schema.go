package database

import (
	"database/sql"
)

type migration struct {
	version int
	script  string
}

var migrations = []migration{
	{
		version: 1,
		script: `
			CREATE TABLE IF NOT EXISTS payments (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				invoice_id TEXT NOT NULL,
				amount NUMERIC NOT NULL,
				date TIMESTAMP NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);

			CREATE TABLE IF NOT EXISTS invoices (
				id TEXT PRIMARY KEY,
				customer_name TEXT NOT NULL,
				amount_due NUMERIC NOT NULL,
				payment_mean TEXT NOT NULL,
				invoice_date TIMESTAMP NOT NULL,
				due_date TIMESTAMP,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
			);
			`,
	},
}

func ApplySchema(db *sql.DB) error {

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`); err != nil {
		return err
	}

	var current int
	row := db.QueryRow(`SELECT IFNULL(MAX(version), 0) FROM schema_migrations`)
	if err := row.Scan(&current); err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version > current {
			if _, err := db.Exec(m.script); err != nil {
				return err
			}
			if _, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, m.version); err != nil {
				return err
			}
		}
	}

	return nil
}
