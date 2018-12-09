package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

// GetAll returns list menu categories sorted by order.
func (MenuCategory) GetAll(db *sqlx.DB) (*[]MenuCategory, error) {
	categories := make([]MenuCategory, 0)

	if err := db.Select(&categories, "SELECT * FROM categories ORDER BY `order` ASC;"); err != nil {
		return nil, err
	}

	return &categories, nil
}

// Find returns MenuItem object with specified ID.
func (MenuCategory) Find(db *sqlx.DB, id uint64) (*MenuCategory, error) {
	category := MenuCategory{}

	if err := db.Get(&category, "SELECT * FROM categories WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &category, nil
}

// Insert adds new category.
func (menuCategory *MenuCategory) Insert(db *sqlx.DB) error {
	sqlStatement := "INSERT INTO categories (name, `order`) VALUES (?, ?);"

	result, err := db.Exec(sqlStatement, menuCategory.Name, menuCategory.Order)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	createdCategory, err := MenuCategory.Find(MenuCategory{}, db, uint64(id))
	if err != nil {
		return err
	}
	*menuCategory = *createdCategory

	return nil
}

// Update menu category object in DB.
func (menuCategory *MenuCategory) Update(db *sqlx.DB) error {

	if _, err := MenuCategory.Find(MenuCategory{}, db, menuCategory.ID); err != nil {
		return err
	}

	// Order validation
	existingCategory := MenuCategory{}
	err := db.Get(&existingCategory, "SELECT * FROM categories WHERE `order` = ?", menuCategory.Order)
	if err == nil && existingCategory.ID != menuCategory.ID {
		return errors.New("item with such ordinal number already exists, change order")
	}

	query := "UPDATE categories SET name=:name, `order`=:order WHERE id=:id"
	_, err = db.NamedExec(query, menuCategory)

	if err != nil {
		return err
	}

	return nil
}
