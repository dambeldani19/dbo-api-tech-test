package entity

import (
	"database/sql"
	"time"
)

type Order struct {
	ID          int
	OrderCode   string
	CustomerID  int
	TotalAmount int
	Status      string
	OrderDetail []OrderDetail `gorm:"foreignkey:OrderID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

type OrderDetail struct {
	ID         int
	OrderID    int
	ProductID  int
	Product    Product `gorm:"foreignkey:ProductID"`
	Quantity   int
	Price      int
	TotalPrice int
}
