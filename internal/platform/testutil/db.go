package testutil

import (
	"app/internal/platform/database"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func setupInMemoryDB() *sql.DB {

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}

	if err := database.ApplySchema(db); err != nil {
		panic(err)
	}

	return db

}

func SetupTestDB(t *testing.T) *sql.DB {
	db := setupInMemoryDB() // your helper
	t.Cleanup(func() {
		db.Close()
	})
	return db
}
