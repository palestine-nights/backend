package db

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

// Initialize DB.
func Initialize(connectionString string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionString)
	db.AutoMigrate(&Table{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// GetAll returns list of all tables.
func (Table) GetAll(db *gorm.DB) *[]Table {
	var tables []Table

	db.Find(&tables)

	return &tables
}

// Find returns Tables object with specified ID.
func (Table) Find(db *gorm.DB, id int) *Table {
	table := Table{}

	db.First(&table, id)

	return &table
}

// Destroy table with specified ID.
func (Table) Destroy(db *gorm.DB, id int) error {
	tbl := Table.Find(Table{}, db, id)

	if tbl.ID == 0 {
		return errors.New("table with such id not found")
	}

	db.Delete(Table{}, id)
	db.Delete(Table{}, "email LIKE ?", "%jinzhu%")

	return nil
}

// Save table object to DB.
func (table *Table) Save(db *gorm.DB) {
	db.Save(table)
}

// Update table object in DB.
func (table *Table) Update(db *gorm.DB) error {
	tbl := Table.Find(Table{}, db, table.ID)

	if tbl.ID == 0 {
		return errors.New("table with such id not found")
	}

	db.Save(table)
	return nil
}
