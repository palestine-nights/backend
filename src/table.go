package main

import (
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

func (table) TableName() string {
	return "tables"
}