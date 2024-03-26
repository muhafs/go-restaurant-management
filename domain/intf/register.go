package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
)

type RegisterUsecase interface {
	Create(c context.Context, user *entity.User) error

	FindByEmail(c context.Context, email string) (entity.User, error)

	FindByPhone(c context.Context, phone string) (entity.User, error)

	CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error)

	CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error)
}
