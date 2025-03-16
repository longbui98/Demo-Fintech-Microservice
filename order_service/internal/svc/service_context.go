package svc

import (
	"micro-project/order_service/internal/config"
	database "micro-project/order_service/internal/db"

	"gorm.io/gorm"
)

type ServiceContext struct {
	DB *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		DB: database.InitDB(c.Database.DataSource),
	}
}
