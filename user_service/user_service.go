package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	CreatedBy string
}

type UserService struct {
	db *gorm.DB
}

type UserArgs struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type UserLoginArgs struct {
	Email    string
	Password string
}

type UserResponse struct {
	Message string
	UserID  uint
}

func initDB() (*gorm.DB, error) {
	dsn := "root:tonyking123@tcp(localhost:3306)/UserDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *UserService) Register(args *UserArgs, response *UserResponse) (*UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := User{
		FirstName: args.FirstName,
		LastName:  args.LastName,
		Email:     args.Email,
		Password:  string(hashedPassword),
	}

	tx := s.db.Begin()
	// Save the user to the database
	if err := s.db.Create(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	response.Message = "User registered successfully"
	response.UserID = user.ID
	return response, nil
}

func (s *UserService) Login(args *UserLoginArgs, response *UserResponse) error {
	var user User

	if err := s.db.Where("email = ?", args.Email).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(args.Password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	response.Message = "Login successful"
	response.UserID = user.ID
	return nil
}

func startServer() {
	db, err := initDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	userService := &UserService{db: db}

	if err := rpc.Register(userService); err != nil {
		log.Fatalf("Error registering User Service: %v", err)
	}

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8000...")

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
