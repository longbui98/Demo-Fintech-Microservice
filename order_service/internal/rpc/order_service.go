package rpc

import (
	"context"
	"micro-project/order_service/internal/logic"
	"micro-project/order_service/internal/model"
	"micro-project/order_service/internal/svc"
)

type OrderService struct {
	ctx *svc.ServiceContext
}

func NewOrderService(ctx *svc.ServiceContext) *OrderService {
	return &OrderService{ctx: ctx}
}

func (s *OrderService) CreateOrder(ctx context.Context, order *model.OrderArgs, res *string) error {
	orderlogic := logic.NewOrderLogic(ctx, s.ctx)
	return orderlogic.CreateOrder(order, res)
}

func (s *OrderService) CreateOrders(ctx context.Context, req *[]model.OrderArgs, res *string) error {
	orderlogic := logic.NewOrderLogic(ctx, s.ctx)
	err := orderlogic.CreateOrders(req, res)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetOrders(ctx context.Context, req *model.Email, res *[]model.OrderResponse) error {
	orderlogic := logic.NewOrderLogic(ctx, s.ctx)
	err := orderlogic.GetOrderByEmail(req, res)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) DeleteOrders(ctx context.Context, req *model.OrderId, res *string) error {
	orderlogic := logic.NewOrderLogic(ctx, s.ctx)
	err := orderlogic.DeleteOrderByOrderId(req, res)
	if err != nil {
		return err
	}
	return nil
}
