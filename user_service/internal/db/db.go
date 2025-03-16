package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(DataSource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	log.Println("connected successfully")

	return db
}
