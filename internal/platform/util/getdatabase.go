package util

import (
	"app/internal/platform/config"
	"app/internal/platform/database"
	"database/sql"
	"log"
)

func GetDB() *sql.DB {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return db

}

func GetDBWithSchema() *sql.DB {
	db := GetDB()

	if err := database.ApplySchema(db); err != nil {
		log.Fatal(err)
		panic(err)
	}

	return db
}
