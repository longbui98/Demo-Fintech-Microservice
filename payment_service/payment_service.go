package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Payment struct {
	PaymentID uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"unique"`
	Method    string
	Status    string
}

type PaymentService struct {
	db *gorm.DB
}

type PaymentArgs struct {
	OrderId uint
	Method  string
}

func initDB() (*gorm.DB, error) {
	dsn := "root:tonyking123@tcp(localhost:3306)/UserDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Payment{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *PaymentService) ProcessPayment(args *PaymentArgs) error {
	if args.OrderId == 0 {
		log.Fatalf("Invalid Order Id")
		return errors.New("invalid Order Id")
	}

	newPayment := Payment{
		OrderID: args.OrderId,
		Method:  args.Method,
		Status:  "Paid",
	}

	result := s.db.Create(&newPayment)
	if result.Error != nil {
		return result.Error
	}

	fmt.Printf("Payment processed for OrderID: %d using %s", args.OrderId, args.Method)
	return nil
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	userService := &PaymentService{db: db}

	if err := rpc.Register(userService); err != nil {
		log.Fatalf("Error registering Payment Service: %v", err)
	}

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8002...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
