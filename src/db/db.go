package db

import "github.com/jmoiron/sqlx"

// Initialize DB.
func Initialize(connectionString string) *sqlx.DB {
	db, err := sqlx.Open("mysql", connectionString)

	if err != nil {
		panic(err)
	}

	return db
}
