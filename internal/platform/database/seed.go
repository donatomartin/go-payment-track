package database

import (
	sql "database/sql"
)

func GetInsertSampleDataQuery() string {
	return `
		-- Delete all existing data
		DELETE FROM payments;

		-- Insert sample payment data
		INSERT INTO payments (invoice_id, amount, date) VALUES
			('1', 100.00, '2023-01-01'),
			('2', 200.00, '2023-01-02'),
			('3', 300.00, '2023-01-03'),
			('4', 400.00, '2023-01-04'),
			('5', 500.00, '2023-01-05'),
			('6', 600.00, '2023-01-06'),
			('7', 700.00, '2023-01-07');
		`
}

func InsertSampleData(db *sql.DB) error {

	_, err := db.Exec(GetInsertSampleDataQuery())
	if err != nil {
		return err
	}

	return nil

}
