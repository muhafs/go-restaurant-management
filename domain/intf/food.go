package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
)

type FoodRepository interface {
	Create(c context.Context, request *entity.Food) (err error)

	Find(c context.Context) (foods []entity.Food, err error)

	FindOne(c context.Context, id string) (food entity.Food, err error)

	FindMenu(c context.Context, id string) (menu entity.Menu, err error)

	Update(c context.Context, request *entity.Food, id string) (err error)

	Delete(c context.Context, id string) (err error)
}

type FoodUsecase interface {
	Create(c context.Context, request *entity.Food) (err error)

	Find(c context.Context) (foods []entity.Food, err error)

	FindOne(c context.Context, id string) (food entity.Food, err error)

	FindMenu(c context.Context, id string) (menu entity.Menu, err error)

	Update(c context.Context, request *entity.Food, id string) (err error)

	Delete(c context.Context, id string) (err error)
}
