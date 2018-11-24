package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Initialize DB.
func Initialize(connectionString string) *sqlx.DB {
	db, err := sqlx.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
