package db

import (
	"github.com/jmoiron/sqlx"
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Find returns Token object with specified id.
func (Token) Find(db *sqlx.DB, id uint64) (*Token, error) {
	token := Token{}

	if err := db.Get(&token, "SELECT * FROM tokens WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &token, nil
}

// FindByCode returns Token object with specified code.
func (Token) FindByCode(db *sqlx.DB, code string) (*Token, error) {
	token := Token{}

	if err := db.Get(&token, "SELECT * FROM tokens WHERE code = ?", code); err != nil {
		return nil, err
	}

	return &token, nil
}

// UpdateState updates token's state.
func (token *Token) UpdateState(db *sqlx.DB) error {

	if _, err := Token.Find(Token{}, db, token.ID); err != nil {
		return err
	}

	_, err := db.NamedExec(`UPDATE tokens SET state=:state WHERE id = :id`, token)

	if err != nil {
		return err
	}

	return nil
}

// Insert adds new token.
func (token *Token) Insert(db *sqlx.DB) error {
	sqlStatement := `INSERT INTO tokens (reservation_id, code, type, state) VALUES (?, ?, ?, ?);`

	result, err := db.Exec(sqlStatement, token.ReservationID, token.Code, token.Type, token.State)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	createdToken, err := Token.Find(Token{}, db, uint64(id))
	if err != nil {
		return err
	}
	*token = *createdToken

	return nil
}

// Generate generates new token.
func (Token) Generate(reservationID uint64, tokenType string) *Token {
	rand.Seed(time.Now().UnixNano())
	code := make([]rune, 64)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	token := Token{}
	token.ReservationID = reservationID
	token.State = "unused"
	token.Type = tokenType
	token.Code = string(code)
	return &token
}
