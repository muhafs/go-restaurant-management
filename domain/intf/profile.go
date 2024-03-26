package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/model"
)

type ProfileUsecase interface {
	Find(c context.Context) ([]*model.ProfileResponse, error)
	FindOne(c context.Context, userID string) (*model.ProfileResponse, error)
}
