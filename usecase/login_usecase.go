package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
)

type LoginUsecase struct {
	UserRepository intf.UserRepository
	ContextTimeout time.Duration
}

func NewLoginUsecase(userRepository intf.UserRepository, timeout time.Duration) intf.LoginUsecase {
	return &LoginUsecase{
		UserRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (lu *LoginUsecase) FindByEmail(c context.Context, email string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(c, lu.ContextTimeout)
	defer cancel()

	user, err = lu.UserRepository.FindByEmail(ctx, email)

	return
}

func (lu *LoginUsecase) CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	accessToken, err = helpers.CreateAccessToken(user, secret, expiry)

	return
}

func (lu *LoginUsecase) CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error) {
	refreshToken, err = helpers.CreateRefreshToken(user, secret, expiry)

	return
}
