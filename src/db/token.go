package db

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
)

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

	_, err := db.NamedExec(`UPDATE tokens SET used=:used WHERE id = :id`, token)

	if err != nil {
		return err
	}

	return nil
}

// Insert adds new token.
func (token *Token) Insert(db *sqlx.DB) error {
	sqlStatement := `INSERT INTO tokens (reservation_id, code, type) VALUES (?, ?, ?);`

	result, err := db.Exec(sqlStatement, token.ReservationID, token.Code, token.Type)

	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	createdToken, err := token.Find(db, uint64(id))
	if err != nil {
		return err
	}
	*token = *createdToken

	return nil
}

// GenerateToken generates new confirmation or cancellation token.
func GenerateToken(reservationID uint64, tokenType TokenType, db *sqlx.DB) error {
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	code := make([]rune, 16)

	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}

	token := Token{
		ReservationID: reservationID,
		Code:          string(code),
		Type:          tokenType,
	}

	err := token.Insert(db)

	if err != nil {
		return err
	}

	return nil
}

// LastToken returns token last inserted token with specified type.
// Returns error in case of SQL error.
func LastToken(db *sqlx.DB, tokenType TokenType) (*Token, error) {
	token := Token{}
	sql := `SELECT * FROM tokens WHERE type = ? ORDER BY ID DESC LIMIT 1;`

	if err := db.Get(&token, sql, tokenType); err != nil {
		return nil, err
	}

	return &token, nil
}

// MustLastToken returns token last inserted token with specified type.
func MustLastToken(db *sqlx.DB, tokenType TokenType) *Token {
	lastToken, err := LastToken(db, tokenType)

	if err != nil {
		panic(err)
	}

	return lastToken
}
