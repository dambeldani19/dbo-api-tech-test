package entity

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}
