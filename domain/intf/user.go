package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
)

type UserRepository interface {
	Create(c context.Context, user *entity.User) error

	Find(c context.Context) ([]entity.User, error)

	FindOne(c context.Context, id string) (entity.User, error)

	FindByEmail(c context.Context, email string) (entity.User, error)
}
