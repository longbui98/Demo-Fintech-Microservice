package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	OrderID   uint `gorm:"primaryKey"`
	UserID    uint
	Name      string
	Quality   int
	Price     int
	Status    OrderStatus
	CreatedAt time.Time
	CreatedBy string
}

type OrderStatus int

const (
	Pending OrderStatus = iota
	Paid
)

type OrderService struct {
	db *gorm.DB
}

type OrderArgs struct {
	UserID  uint
	Name    string
	Quality int
	Price   int
}

type PurschaseOrderArgs struct {
}

type OrderResponse struct {
	OrderId uint
	UserId  uint
	Name    string
	Quality int
	Status  int
}

func initDB() (*gorm.DB, error) {
	dsn := "root:tonyking123@tcp(localhost:3306)/UserDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Order{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *OrderService) CreateOrder(args *OrderArgs) error {
	order := Order{
		UserID:  args.UserID,
		Name:    args.Name,
		Quality: args.Quality,
		Price:   args.Price,
		Status:  Pending,
	}
	tx := s.db.Begin()

	if err := s.db.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	fmt.Printf("created order successully")
	return nil
}

func (s *OrderService) CreateOders(args *[]OrderArgs) error {
	var orders []Order
	for _, order := range *args {
		orders = append(orders,
			Order{
				UserID:  order.UserID,
				Quality: order.Quality,
				Name:    order.Name,
				Price:   order.Price,
				Status:  Pending,
			})
	}

	tx := s.db.Begin()
	if err := s.db.CreateInBatches(&orders, 100).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	fmt.Printf("created orders successully")

	return nil
}

func (s *OrderService) GetOrderByUserId(user_id uint, res *OrderResponse) ([]*OrderResponse, error) {
	var order []Order
	result := s.db.Where("user_id = ?", user_id).Find(&order)

	if result.Error != nil {
		return nil, result.Error
	}

	var responseOrders []*OrderResponse
	for _, order := range order {
		responseOrders = append(responseOrders,
			&OrderResponse{
				OrderId: order.OrderID,
				UserId:  order.UserID,
				Name:    order.Name,
				Quality: order.Quality,
				Status:  int(order.Status),
			})
	}
	return responseOrders, nil
}

func (s *OrderService) DeleteOrderByOrderId(order_id []uint) error {
	if len(order_id) == 0 {
		return fmt.Errorf("list of order_id is empty")
	}

	tx := s.db.Begin()
	result := s.db.Where("order_id IN ?", order_id).Delete(&Order{})

	if result.RowsAffected == 0 || result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("cannot delete order_id")
	}

	tx.Commit()
	fmt.Printf("deleted list of Ids successully")

	return nil
}

func startServer() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	orderService := &OrderService{db: db}

	if err := rpc.Register(orderService); err != nil {
		log.Fatalf("Error registering Order Service: %v", err)
	}

	listener, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8001...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func main() {
	startServer()
}
