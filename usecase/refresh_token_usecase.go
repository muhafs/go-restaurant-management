package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
)

type RefreshTokenUsecase struct {
	UserRepository intf.UserRepository
	ContextTimeout time.Duration
}

func NewRefreshTokenUsecase(userRepository intf.UserRepository, timeout time.Duration) intf.RefreshTokenUsecase {
	return &RefreshTokenUsecase{
		UserRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (ru *RefreshTokenUsecase) FindOne(c context.Context, email string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(c, ru.ContextTimeout)
	defer cancel()

	user, err = ru.UserRepository.FindOne(ctx, email)

	return
}

func (ru *RefreshTokenUsecase) CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	accessToken, err = helpers.CreateAccessToken(user, secret, expiry)

	return
}

func (ru *RefreshTokenUsecase) CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error) {
	refreshToken, err = helpers.CreateRefreshToken(user, secret, expiry)

	return
}

func (ru *RefreshTokenUsecase) ExtractIDFromToken(requestToken string, secret string) (id string, err error) {
	id, err = helpers.ExtractIDFromToken(requestToken, secret)

	return
}
