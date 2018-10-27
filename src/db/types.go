package db

// Table object for REST API.
type Table struct {
	ID          int    `json:"id" gorm:"primary_key"`
	Places      int    `json:"places"`
	Description string `json:"description"`
}
