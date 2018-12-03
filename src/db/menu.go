package db

import (
	"github.com/jmoiron/sqlx"
)

// GetAll returns list of all menu items.
func (MenuItem) GetAll(db *sqlx.DB) (*[]MenuItem, error) {
	menuItems := make([]MenuItem, 0)

	if err := db.Select(&menuItems, `SELECT * FROM menu;`); err != nil {
		return nil, err
	}

	return &menuItems, nil
}

// GetByCategory returns list of all menu items in given category.
func (MenuItem) GetByCategory(db *sqlx.DB, category string) (*[]MenuItem, error) {
	menuItems := make([]MenuItem, 0)

	if err := db.Select(&menuItems, `SELECT * FROM menu WHERE category = ?`, category); err != nil {
		return nil, err
	}

	return &menuItems, nil
}

// Find returns MenuItem object with specified ID.
func (MenuItem) Find(db *sqlx.DB, id uint64) (*MenuItem, error) {
	menuItem := MenuItem{}

	if err := db.Get(&menuItem, "SELECT * FROM menu WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &menuItem, nil
}

// Destroy menu item with specified ID.
func (MenuItem) Destroy(db *sqlx.DB, id uint64) error {

	if _, err := MenuItem.Find(MenuItem{}, db, id); err != nil {
		return err
	}

	if _, err := db.Exec(`DELETE FROM menu WHERE id = ?;`, id); err != nil {
		return err
	}

	return nil
}

// Update menu item object in DB.
func (menuItem *MenuItem) Update(db *sqlx.DB) error {

	if _, err := MenuItem.Find(MenuItem{}, db, menuItem.ID); err != nil {
		return err
	}

	query := `UPDATE menu SET name=:name, description=:description, price=:price, category=:category, image_url=:image_url WHERE id=:id`
	_, err := db.NamedExec(query, menuItem)

	if err != nil {
		return err
	}

	return nil
}

// Insert adds new menu item.
func (menuItem *MenuItem) Insert(db *sqlx.DB) error {
	sqlStatement := `INSERT INTO menu (name, description, price, category, image_url) VALUES (?, ?, ?, ?, ?);`

	result, err := db.Exec(sqlStatement, menuItem.Name, menuItem.Description, menuItem.Price, menuItem.Category, menuItem.ImageURL)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	createdItem, err := MenuItem.Find(MenuItem{}, db, uint64(id))
	if err != nil {
		return err
	}
	*menuItem = *createdItem

	return nil
}

// GetCategories returns list of unique menu categories.
func (MenuItem) GetCategories(db *sqlx.DB) ([]string, error) {
	categories := make([]string, 0)

	if err := db.Select(&categories, `SELECT category FROM menu GROUP by category;`); err != nil {
		return nil, err
	}

	return categories, nil
}
