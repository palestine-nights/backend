package db

import "time"

// Table object for REST API.
type Table struct {
	ID          uint64    `json:"id" db:"id"`
	Places      int64     `json:"places" db:"places"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// Reservation object for REST API.
type Reservation struct {
	ID        uint64        `json:"id" db:"id"`
	TableID   uint64        `json:"table_id" db:"table_id"`
	Guests    int64         `json:"guests" db:"guests"`
	Email     string        `json:"email" db:"email"`
	Phone     string        `json:"phone" db:"phone"`
	Status    string        `json:"status" db:"status"`
	FullName  string        `json:"fullname" db:"fullname"`
	Time      time.Time     `json:"time" db:"time"`
	Duration  time.Duration `json:"duration" db:"duration"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
}
