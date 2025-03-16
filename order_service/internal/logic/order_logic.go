package logic

import (
	"context"
	"fmt"
	"micro-project/order_service/internal/model"
	"micro-project/order_service/internal/svc"

	"github.com/go-playground/validator/v10"
)

type OrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderLogic {
	return &OrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var validate = validator.New()

func ValidateStruct(order model.OrderArgs) error {
	return validate.Struct(order)
}

func (s *OrderLogic) CreateOrder(args *model.OrderArgs, response *string) error {
	if err := ValidateStruct(*args); err != nil {
		return err
	}
	order := model.Order{
		Email:    args.Email,
		Name:     args.Name,
		Quantity: args.Quantity,
		Price:    args.Price,
		Status:   model.Pending,
	}
	tx := s.svcCtx.DB.Begin()

	if err := s.svcCtx.DB.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	fmt.Print("Result created order: ", order)

	*response = "Created successully"

	tx.Commit()
	fmt.Printf("created order successully")
	return nil
}

func (s *OrderLogic) CreateOrders(args *[]model.OrderArgs, response *string) error {
	for _, val := range *args {
		if err := ValidateStruct(val); err != nil {
			return err
		}
	}

	var orders []model.Order
	for _, order := range *args {
		orders = append(orders,
			model.Order{
				Email:    order.Email,
				Quantity: order.Quantity,
				Name:     order.Name,
				Price:    order.Price,
				Status:   model.Pending,
			})
	}

	tx := s.svcCtx.DB.Begin()
	if err := s.svcCtx.DB.CreateInBatches(&orders, 100).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	fmt.Printf("created orders successully")

	*response = "Created successully"

	return nil
}

func (s *OrderLogic) GetOrderByEmail(email *model.Email, res *[]model.OrderResponse) error {
	var order []model.Order
	result := s.svcCtx.DB.Where("email = ?", email.Email).Find(&order)

	if result.Error != nil {
		fmt.Errorf("error while get or by email: ", result.Error)
		return nil
	}

	var responseOrders []model.OrderResponse
	for _, order := range order {
		responseOrders = append(responseOrders,
			model.OrderResponse{
				Email:    order.Email,
				Name:     order.Name,
				Quantity: order.Quantity,
				Status:   int(order.Status),
			})
	}

	fmt.Println(responseOrders)

	*res = responseOrders
	return nil
}

func (s *OrderLogic) DeleteOrderByOrderId(order_id *model.OrderId, response *string) error {
	if order_id == nil {
		return fmt.Errorf("order_id cannot empty")
	}

	tx := s.svcCtx.DB.Begin()
	result := s.svcCtx.DB.Where("order_id = ?", order_id.OrderID).Delete(&model.Order{})

	if result.RowsAffected == 0 || result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("cannot delete order_id")
	}

	*response = "Deleted successully"

	tx.Commit()
	fmt.Printf("deleted list of Ids successully")

	return nil
}
