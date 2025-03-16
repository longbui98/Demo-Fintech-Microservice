package model

import (
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
