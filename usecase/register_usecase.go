package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
)

type RegisterUsecase struct {
	UserRepository intf.UserRepository
	ContextTimeout time.Duration
}

func NewRegisterUsecase(userRepository intf.UserRepository, timeout time.Duration) intf.RegisterUsecase {
	return &RegisterUsecase{
		UserRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (ru *RegisterUsecase) Create(c context.Context, user *entity.User) (err error) {
	ctx, cancel := context.WithTimeout(c, ru.ContextTimeout)
	defer cancel()

	err = ru.UserRepository.Create(ctx, user)

	return
}

func (ru *RegisterUsecase) FindByEmail(c context.Context, email string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(c, ru.ContextTimeout)
	defer cancel()

	user, err = ru.UserRepository.FindByEmail(ctx, email)

	return
}

func (ru *RegisterUsecase) FindByPhone(c context.Context, phone string) (user entity.User, err error) {
	ctx, cancel := context.WithTimeout(c, ru.ContextTimeout)
	defer cancel()

	user, err = ru.UserRepository.FindByPhone(ctx, phone)

	return
}

func (ru *RegisterUsecase) CreateAccessToken(user *entity.User, secret string, expiry int) (accessToken string, err error) {
	accessToken, err = helpers.CreateAccessToken(user, secret, expiry)

	return
}

func (ru *RegisterUsecase) CreateRefreshToken(user *entity.User, secret string, expiry int) (refreshToken string, err error) {
	refreshToken, err = helpers.CreateRefreshToken(user, secret, expiry)

	return
}
