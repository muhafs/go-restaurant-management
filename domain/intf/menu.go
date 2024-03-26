package intf

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
)

type MenuRepository interface {
	Create(c context.Context, request *entity.Menu) (err error)

	Find(c context.Context) (menus []entity.Menu, err error)

	FindOne(c context.Context, id string) (menu entity.Menu, err error)

	Update(c context.Context, request *entity.Menu, id string) (err error)

	Delete(c context.Context, id string) (err error)
}

type MenuUsecase interface {
	Create(c context.Context, request *entity.Menu) (err error)

	Find(c context.Context) (menus []entity.Menu, err error)

	FindOne(c context.Context, id string) (menu entity.Menu, err error)

	Update(c context.Context, request *entity.Menu, id string) (err error)

	Delete(c context.Context, id string) (err error)
}
