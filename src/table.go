package main

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type table struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Places      int    `json:"places"`
	Description string `json:"description"`
}

func (table) getTable(db *gorm.DB, id int) *table {
	t := table{}
	db.First(&t, id)
	return &t
}

func (t *table) createTable(db *gorm.DB) {
	db.Save(t)
}

func (t *table) updateTable(db *gorm.DB) error {
	tbl := table.getTable(table{}, db, t.ID)
	if tbl.ID == 0 {
		return errors.New("table with such id not found")
	}
	db.Save(t)
	return nil
}

func (table) deleteTable(db *gorm.DB, id int) error {
	tbl := table.getTable(table{}, db, id)
	if tbl.ID == 0 {
		return errors.New("table with such id not found")
	}
	db.Delete(table{}, id)
	return nil
}

func (table) TableName() string {
	return "tables"
}