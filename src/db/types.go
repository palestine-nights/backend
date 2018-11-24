package db

import (
	"time"
)

// Table object for REST API.
type Table struct {
	ID          uint64    `json:"id" db:"id"`
	Places      int64     `json:"places" db:"places"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// State is table reservation state.
// Three possible states are availabe: 'created', 'approved', 'cancelled'.
type State string

const (
	// StateCreated returns state string of created reservation.
	StateCreated State = "created"
	// StateApproved returns state string of approved reservation.
	StateApproved State = "approved"
	// StateCancelled returns state string of cancelled reservation.
	StateCancelled State = "cancelled"
)

// Reservation object for REST API.
type Reservation struct {
	ID        uint64        `json:"id" db:"id"`
	TableID   uint64        `json:"table_id" db:"table_id"`
	Guests    int64         `json:"guests" db:"guests"`
	Email     string        `json:"email" db:"email"`
	Phone     string        `json:"phone" db:"phone"`
	State     State         `json:"state" db:"state"`
	FullName  string        `json:"full_name" db:"full_name"`
	Time      time.Time     `json:"time" db:"time"`
	Duration  time.Duration `json:"duration" db:"duration"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
}
