package db

import (
	"github.com/jmoiron/sqlx"
)

// GetAll returns list of all tables.
func (Table) GetAll(db *sqlx.DB) (*[]Table, error) {
	tables := make([]Table, 0)

	if err := db.Select(&tables, `SELECT * FROM tables;`); err != nil {
		return nil, err
	}

	return &tables, nil
}

// Find returns Tables object with specified ID.
func (Table) Find(db *sqlx.DB, id uint64) (*Table, error) {
	table := Table{}

	if err := db.Get(&table, "SELECT * FROM tables WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &table, nil
}

// Destroy table with specified ID.
func (Table) Destroy(db *sqlx.DB, id uint64) error {

	if _, err := Table.Find(Table{}, db, id); err != nil {
		return err
	}

	if _, err := db.Exec(`DELETE FROM tables WHERE id = ?;`, id); err != nil {
		return err
	}

	return nil
}

// Update table object in DB.
func (table *Table) Update(db *sqlx.DB) error {

	if _, err := Table.Find(Table{}, db, table.ID); err != nil {
		return err
	}

	query := `UPDATE tables SET places=:places, description=:description, active=:active WHERE id = :id`
	_, err := db.NamedExec(query, table)

	if err != nil {
		return err
	}

	return nil
}

// Insert adds new table.
func (table *Table) Insert(db *sqlx.DB) error {
	sqlStatement := `INSERT INTO tables (places, description, active) VALUES (?, ?, ?);`

	result, err := db.Exec(sqlStatement, table.Places, table.Description, table.Active)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	createdTable, err := Table.Find(Table{}, db, uint64(id))
	if err != nil {
		return err
	}
	*table = *createdTable

	return nil
}
