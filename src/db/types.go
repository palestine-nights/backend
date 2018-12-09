package db

import "time"

// Table represents restaurant table.
//
// swagger:parameters model
type Table struct {
	ID uint64 `json:"id" db:"id"`
	// Number of places to seat.
	// required: true
	Places int64 `json:"places" db:"places"`
	// Description of the table.
	// required: true
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// State is string representation of reservation state.
// swagger:strfmt state
type State string

const (
	// StateCreated returns state string of created reservation.
	StateCreated State = "created"
	// StateApproved returns state string of approved reservation.
	StateApproved State = "approved"
	// StateCancelled returns state string of cancelled reservation.
	StateCancelled State = "cancelled"
)

// Reservation model for table reservation proccess.
//
// swagger:model
type Reservation struct {
	ID uint64 `json:"id" db:"id"`
	// ID of table, associated with reservation.
	// required: true
	TableID uint64 `json:"table_id" db:"table_id"`
	// Number of people to seat for reservation.
	// required: true
	Guests int64 `json:"guests" db:"guests"`
	// Email of the client.
	// required: true
	Email string `json:"email" db:"email"`
	// Phone of the client.
	// required: true
	Phone string `json:"phone" db:"phone"`
	State State  `json:"state" db:"state"`
	// Full Name of the client.
	// required: true
	FullName string `json:"full_name" db:"full_name"`
	// Time of the reservation.
	// required: true
	Time time.Time `json:"time" db:"time"`
	// Duration of the reservation.
	// required: truee
	Duration  time.Duration `json:"duration" db:"duration"`
	CreatedAt time.Time     `json:"-" db:"created_at"`
	UpdatedAt time.Time     `json:"-" db:"updated_at"`
}

// MenuItem model for menu.
//
// swagger:model
type MenuItem struct {
	ID uint64 `json:"id" db:"id"`
	// Name of the menu item.
	// required: true
	Name string `json:"name" db:"name"`
	// Description of the menu item.
	// required: true
	Description string `json:"description" db:"description"`
	// Price of the menu item in Bahrain Dinars.
	// required: true
	Price float32 `json:"price" db:"price"`
	// Category of the menu item.
	// required: true
	CategoryID uint64 `json:"category_id" db:"category_id"`
	// Image URL for the menu item.
	// required: true
	ImageURL  string    `json:"image_url" db:"image_url"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// MenuCategory model for menu categories.
//
// swagger:model
type MenuCategory struct {
	ID uint64 `json:"id" db:"id"`
	// Name of the menu category.
	// required: true
	Name string `json:"name" db:"name"`
	// Order of this category in categories list.
	// required: true
	Order     uint64    `json:"order" db:"order"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// Token object for REST API.
//
// swagger:model
type Token struct {
	ID            uint64    `json:"id" db:"id"`
	ReservationID uint64    `json:"reservation_id" db:"reservation_id"`
	Code          string    `json:"code" db:"code"`
	Type          string    `json:"type" db:"type"`
	State         string    `json:"state" db:"state"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"-" db:"updated_at"`
}
