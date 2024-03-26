package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
)

type MenuUsecase struct {
	MenuRepository intf.MenuRepository
	ContextTimeout time.Duration
}

func NewMenuUsecase(menuRepository intf.MenuRepository, timeout time.Duration) intf.MenuUsecase {
	return &MenuUsecase{
		MenuRepository: menuRepository,
		ContextTimeout: timeout,
	}
}

// Create implements intf.MenuUsecase.
func (u *MenuUsecase) Create(c context.Context, request *entity.Menu) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.MenuRepository.Create(ctx, request)

	return
}

// Find implements intf.MenuUsecase.
func (u *MenuUsecase) Find(c context.Context) (menus []entity.Menu, err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	menus, err = u.MenuRepository.Find(ctx)

	return
}

// FindOne implements intf.MenuUsecase.
func (u *MenuUsecase) FindOne(c context.Context, id string) (menu entity.Menu, err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	menu, err = u.MenuRepository.FindOne(ctx, id)

	return
}

// Update implements intf.MenuUsecase.
func (u *MenuUsecase) Update(c context.Context, request *entity.Menu, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.MenuRepository.Update(ctx, request, id)

	return
}

// Delete implements intf.MenuUsecase.
func (u *MenuUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.MenuRepository.Delete(ctx, id)

	return
}
