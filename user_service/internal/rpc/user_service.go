package rpc

import (
	"context"
	"micro-project/user_service/internal/logic"
	"micro-project/user_service/internal/model"
	"micro-project/user_service/internal/svc"
)

type UserService struct {
	ctx *svc.ServiceContext
}

func NewUserService(ctx *svc.ServiceContext) *UserService {
	return &UserService{ctx: ctx}
}

func (s *UserService) Register(ctx context.Context, user *model.User, _ *struct{}) error {
	userLogic := logic.NewUserLogic(ctx, s.ctx)
	return userLogic.Register(*user)
}

func (s *UserService) Login(ctx context.Context, req *model.User, res *model.User) error {
	userLogic := logic.NewUserLogic(ctx, s.ctx)
	user, err := userLogic.Login(req.Email, req.Password)
	if err != nil {
		return err
	}

	*res = *user
	return nil
}
