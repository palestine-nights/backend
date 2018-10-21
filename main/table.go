package main

import (
	"database/sql"
	"fmt"
)

type table struct {
	ID			int		`json:"id"`
	Places		int		`json:"places"`
	Description	string	`json:"description"`
}

func (t *table) getTable(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT id, places, description FROM tables WHERE id = ?")
	return db.QueryRow(statement, t.ID).Scan(&t.ID, &t.Places, &t.Description)
}
