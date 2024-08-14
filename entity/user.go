package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
}

type Customer struct {
	ID        int
	UserID    int
	Name      string
	Email     string
	Phone     string
	Address   string
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
