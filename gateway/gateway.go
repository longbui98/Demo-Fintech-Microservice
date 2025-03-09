package main

import (
	"log"
	"net/http"
	"net/rpc"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
	Status   string `json:"status"`
}

type Payment struct {
	OrderID uint   `json:"order_id"`
	Method  string `json:"method"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func UserServiceRegisterClient(user User, response *User) error {
	client, err := rpc.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call("UserService.Register", user, response)
	return err
}

func UserServiceLoginClient(user User, response *User) error {
	client, err := rpc.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call("UserService.Login", user, response)
	return err
}

func OrderServiceClient(order Order, response *Order) error {
	client, err := rpc.Dial("tcp", "localhost:8001")
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call("OrderService.CreateOrder", order, response)
	return err
}

func PurchaseServiceClient(payment Payment, response *string) error {
	client, err := rpc.Dial("tcp", "localhost:8002")
	if err != nil {
		return err
	}
	defer client.Close()

	err = client.Call("PaymentService.ProcessPayment", payment, response)
	return err
}

func main() {
	app := fiber.New()

	app.Post("/orders", func(c *fiber.Ctx) error {
		var order Order
		if err := c.BodyParser(&order); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var response Order
		err := OrderServiceClient(order, &response)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(response)
	})

	app.Post("/payment/:orderId", func(c *fiber.Ctx) error {
		orderID := c.Params("orderId")

		var payment Payment
		if err := c.BodyParser(&payment); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		orderId, errConvert := strconv.Atoi(orderID)
		if errConvert != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Order ID"})
		}

		payment.OrderID = uint(orderId)

		var response string
		err := PurchaseServiceClient(payment, &response)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": response})
	})

	app.Post("/user/register", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var response User
		err := UserServiceRegisterClient(user, &response)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(response)
	})

	app.Post("/user/login", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		var response User
		err := UserServiceLoginClient(user, &response)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(response)
	})

	log.Println("API Gateway running on port 8080...")
	log.Fatal(app.Listen(":8080"))
}
