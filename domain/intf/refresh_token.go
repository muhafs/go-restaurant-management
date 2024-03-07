package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
)

type RefreshTokenUsecase interface {
	FindOne(c context.Context, id string) (entity.User, error)

	CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error)

	CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error)

	ExtractIDFromToken(requestToken string, secret string) (string, error)
}
