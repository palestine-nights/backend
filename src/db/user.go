package db

import "github.com/jmoiron/sqlx"

// FindByPassword returns User object with given username and password.
func (User) FindByPassword(db *sqlx.DB, username, password string) (*User, error) {
	user := User{}

	if err := db.Get(&user, "SELECT * FROM users WHERE username = ? AND password = ?", username, password); err != nil {
		return nil, err
	}

	return &user, nil
}
