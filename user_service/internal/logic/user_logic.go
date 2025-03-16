package logic

import (
	"context"
	"errors"
	"micro-project/user_service/internal/model"
	"micro-project/user_service/internal/svc"

	"golang.org/x/crypto/bcrypt"
)

type UserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) Register(user model.User) error {
	var existingUser model.User
	if err := l.svcCtx.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return l.svcCtx.DB.Create(&user).Error
}

func (l *UserLogic) Login(email, password string) (*model.User, error) {
	var user model.User
	if err := l.svcCtx.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
