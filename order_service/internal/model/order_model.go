package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	OrderID   uint   `gorm:"primaryKey"`
	Email     string `gorm:"not null"`
	Name      string `gorm:"not null"`
	Quantity  int    `gorm:"not null"`
	Price     int    `gorm:"not null"`
	Status    OrderStatus
	CreatedAt time.Time
	CreatedBy string
}

type OrderStatus int

const (
	Pending OrderStatus = iota
	Paid
)

type OrderArgs struct {
	Email    string `gorm:"not null" validate:"required"`
	Name     string `gorm:"not null" validate:"required"`
	Quantity int    `gorm:"not null" validate:"required"`
	Price    int    `gorm:"not null" validate:"required"`
}

type OrderResponse struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Status   int    `json:"status"`
}

type OrderId struct {
	OrderID int
}

type Email struct {
	Email string
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Order{})
}
