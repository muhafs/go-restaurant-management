package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

type ProfileUsecase struct {
	UserRepository intf.UserRepository
	ContextTimeout time.Duration
}

func NewProfileUsecase(userRepository intf.UserRepository, timeout time.Duration) intf.ProfileUsecase {
	return &ProfileUsecase{
		UserRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (pu *ProfileUsecase) Find(c context.Context) (profiles []*model.ProfileResponse, err error) {
	ctx, cancel := context.WithTimeout(c, pu.ContextTimeout)
	defer cancel()

	users, err := pu.UserRepository.Find(ctx)
	if err != nil {
		return
	}

	for _, user := range users {
		profiles = append(profiles, &model.ProfileResponse{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}

	return
}

func (pu *ProfileUsecase) FindOne(c context.Context, userID string) (profile *model.ProfileResponse, err error) {
	ctx, cancel := context.WithTimeout(c, pu.ContextTimeout)
	defer cancel()

	user, err := pu.UserRepository.FindOne(ctx, userID)
	if err != nil {
		return
	}

	profile = &model.ProfileResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return
}
