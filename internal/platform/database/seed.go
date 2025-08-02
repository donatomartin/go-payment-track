package database

import (
	sql "database/sql"
)

func InsertSampleData(db *sql.DB) error {

	insertDataQuery := `
				
				-- Delete all existing data
				DELETE FROM payments;
				
				-- Insert sample payment data
				INSERT INTO payments (amount) VALUES
					(100.00),
					(200.00),
					(300.00),
					(400.00),
					(500.00),
					(600.00)
    `
	_, err := db.Exec(insertDataQuery)
	if err != nil {
		return err
	}

	return nil

}
