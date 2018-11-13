package db

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// Validate ...
func (reservation *Reservation) Validate(db *sqlx.DB) (bool, error) {
	var reservations []Reservation
	var tableReservations []Reservation

	sql := `SELECT * FROM reservations WHERE (created_at >= NOW() - INTERVAL 1 DAY) AND (email != ? OR phone != ?);`
	err := db.Select(reservations, sql, reservation.Email, reservation.Phone)

	if err != nil {
		return false, err
	}

	isValid := len(reservations) == 0

	err = db.Select(tableReservations, `SELECT * FROM reservations WHERE table_id = ?`, reservation.TableID)

	if err != nil {
		return false, err
	}

	// newStartTime := reservation.Time
	// newFinishTime := newStartTime.Add(reservation.Duration)

	return isValid, nil

	// for i, tableReservation := range tableReservations {
	// 	reservationStartTime := tableReservation.Time
	// 	reservationFinishTime := reservationStartTime.Add(tableReservation.Duration)

	// 	if reservationStartTime.Before(newFinishTime) &&
	// 		reservationFinishTime.After(newStartTime) &&
	// 		reservationStartTime.Before(newFinishTime) &&
	// 		newStartTime.Before(reservationFinishTime) {
	// 		return false, nil
	// 	}

	// }

	// return
}

// GetAll returns list of all reservations.
func (Reservation) GetAll(db *sqlx.DB) (*[]Reservation, error) {
	var reservations []Reservation

	err := db.Select(reservations, `SELECT * FROM reservations;`)

	if err != nil {
		return nil, err
	}

	return &reservations, nil
}

// Find returns Reservation's object with specified ID.
func (Reservation) Find(db *sqlx.DB, id uint64) (*Reservation, error) {
	reservation := Reservation{}

	err := db.Get(reservation, `SELECT * FROM tables WHERE id = ?;`)

	if err != nil {
		return nil, errors.New("cannot load reservation with id " + strconv.FormatUint(id, 16))
	}

	return &reservation, nil
}

// Where returns objects.
func (Reservation) Where(db *sqlx.DB, sql string, args ...interface{}) (*[]Reservation, error) {
	var reservations []Reservation

	err := db.Select(reservations, "SELECT * FROM reservations WHERE"+sql)

	if err != nil {
		return nil, errors.New("cannot load reservations")
	}

	return &reservations, nil
}

// Destroy reservation with specified ID.
func (Reservation) Destroy(db *sqlx.DB, id uint64) error {
	_, err := Reservation.Find(Reservation{}, db, id)

	if err != nil {
		panic(err)
	}

	return nil
}
