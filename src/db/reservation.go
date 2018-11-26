package db

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/palestine-nights/backend/src/tools"
	"github.com/ttacon/libphonenumber"
)

// GetStopTime calculates finish time of reservations.
func (reservation *Reservation) GetStopTime() time.Time {
	return reservation.Time.Add(reservation.Duration)
}

func isOverlap(start1, finish1, start2, finish2 time.Time) bool {
	return start1.Before(finish2) && finish1.After(start2) &&
		start1.Before(finish2) && start2.Before(finish1)
}

// Validates time to be not taken to create new table reservation record.
func (reservation *Reservation) validateTime(db *sqlx.DB) error {
	// TODO: List only 'approved' records.
	reservations := make([]Reservation, 0)

	sql := `SELECT * FROM reservations WHERE table_id = ?`

	if err := db.Select(&reservations, sql, reservation.TableID); err != nil {
		return err
	}

	for _, tmp := range reservations {
		if isOverlap(reservation.Time, reservation.GetStopTime(), tmp.Time, tmp.GetStopTime()) {
			return errors.New("This time was already taken")
		}
	}

	return nil
}

// Validate validates all conditions to create new table reservation record.
func (reservation *Reservation) Validate(db *sqlx.DB) (bool, error) {

	if !tools.ValidateEmail(reservation.Email) {
		return false, errors.New("Email is invalid")
	}

	// Validates and formats phone number.
	phoneNumber, err := libphonenumber.Parse(reservation.Phone, "BH")
	if err != nil || !libphonenumber.IsValidNumber(phoneNumber) {
		return false, errors.New("Phone is invalid")
	}
	reservation.Phone = libphonenumber.Format(phoneNumber, libphonenumber.E164)

	reservations := make([]Reservation, 0)

	sql := `SELECT * FROM reservations WHERE (created_at >= NOW() - INTERVAL 1 DAY) AND (email = ? OR phone = ?);`

	if err := db.Select(&reservations, sql, reservation.Email, reservation.Phone); err != nil {
		return false, err
	}

	if len(reservations) != 0 {
		return false, errors.New("Email or phone was already used for last 24 hours")
	}

	if err := reservation.validateTime(db); err != nil {
		return false, err
	}

	return true, nil
}

// GetAll returns list of all reservations.
func (Reservation) GetAll(db *sqlx.DB) (*[]Reservation, error) {
	reservations := make([]Reservation, 0)

	if err := db.Select(&reservations, `SELECT * FROM reservations;`); err != nil {
		return nil, err
	}

	return &reservations, nil
}

// GetUpcoming returns upcoming reservations.
func (Reservation) GetUpcoming(db *sqlx.DB) (*[]Reservation, error) {
	reservations := make([]Reservation, 0)

	if err := db.Select(&reservations, "SELECT * FROM reservations WHERE time >= NOW()"); err != nil {
		return nil, err
	}

	return &reservations, nil
}

// Find returns Reservation's object with specified ID.
func (Reservation) Find(db *sqlx.DB, id uint64) (*Reservation, error) {
	reservation := Reservation{}

	if err := db.Get(&reservation, `SELECT * FROM reservations WHERE id = ?;`, id); err != nil {
		return nil, err
	}

	return &reservation, nil
}

// Destroy reservation with specified ID.
func (Reservation) Destroy(db *sqlx.DB, id uint64) error {
	if _, err := Reservation.Find(Reservation{}, db, id); err != nil {
		return err
	}

	return nil
}

// Insert adds new reservation.
func (reservation *Reservation) Insert(db *sqlx.DB) error {
	sql := `INSERT INTO reservations (table_id,guests,email,phone,full_name,time,duration) VALUES (?, ?, ?, ?, ?, ?, ?)`

	res, err := db.Exec(sql,
		reservation.TableID,
		reservation.Guests,
		reservation.Email,
		reservation.Phone,
		reservation.FullName,
		reservation.Time,
		reservation.Duration,
	)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return err
	}

	createdReservation, err := Reservation.Find(Reservation{}, db, uint64(id))
	*reservation = *createdReservation

	if err != nil {
		return err
	}

	return nil
}
