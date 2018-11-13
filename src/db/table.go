package db

import (
	"errors"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// GetAll returns list of all tables.
func (Table) GetAll(db *sqlx.DB) (*[]Table, error) {
	var tables []Table

	rows, err := db.Queryx(`SELECT * FROM tables;`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table Table

		err = rows.StructScan(&table)

		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	return &tables, nil
}

// Find returns Tables object with specified ID.
func (Table) Find(db *sqlx.DB, id uint64) (*Table, error) {
	table := Table{}

	err := db.Get(&table, "SELECT * FROM tables WHERE id = ?", id)

	if err != nil {
		return nil, errors.New("cannot load table with id " + strconv.FormatUint(id, 16))
	}

	return &table, nil
}

// Destroy table with specified ID.
func (Table) Destroy(db *sqlx.DB, id uint64) error {
	_, err := Table.Find(Table{}, db, id)

	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM tables WHERE id = ?;`, id)

	if err != nil {
		return err
	}

	return nil
}

// Update table object in DB.
func (table *Table) Update(db *sqlx.DB) error {
	_, err := Table.Find(Table{}, db, table.ID)

	if err != nil {
		return err
	}

	_, err = db.NamedExec(`UPDATE tables SET places=:places, description=:description WHERE id = id`, table)

	if err != nil {
		return err
	}

	return nil
}

// Insert adds new table.
func (table *Table) Insert(db *sqlx.DB) error {
	sqlStatement := `INSERT INTO tables (places,description,created_at,updated_at) VALUES (?, ?, ?, ?);`

	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()

	result := db.MustExec(sqlStatement, table.Places, table.Description, time.Now(), time.Now())

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	table.ID = uint64(id)

	return nil
}
