package repository

import (
	"database/sql"
)

func ApplySchema(db *sql.DB) error {

	schema := `
    CREATE TABLE IF NOT EXISTS payments (
        id SERIAL PRIMARY KEY,
        amount NUMERIC(10, 2),
        created_at TIMESTAMP DEFAULT now()
    );
    `
	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	return nil

}
