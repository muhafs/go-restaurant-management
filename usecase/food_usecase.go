package usecase

import (
	"context"
	"time"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
)

type FoodUsecase struct {
	FoodRepository intf.FoodRepository
	ContextTimeout time.Duration
}

func NewFoodUsecase(foodRepository intf.FoodRepository, timeout time.Duration) intf.FoodUsecase {
	return &FoodUsecase{
		FoodRepository: foodRepository,
		ContextTimeout: timeout,
	}
}

// Create implements intf.FoodUsecase.
func (u *FoodUsecase) Create(c context.Context, request *entity.Food) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.FoodRepository.Create(ctx, request)

	return
}

// Find implements intf.FoodUsecase.
func (u *FoodUsecase) Find(c context.Context) (foods []entity.Food, err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	foods, err = u.FoodRepository.Find(ctx)

	return
}

// FindOne implements intf.FoodUsecase.
func (u *FoodUsecase) FindOne(c context.Context, id string) (food entity.Food, err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	food, err = u.FoodRepository.FindOne(ctx, id)

	return
}

// FindOne implements intf.FoodUsecase.
func (u *FoodUsecase) FindMenu(c context.Context, id string) (menu entity.Menu, err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	menu, err = u.FoodRepository.FindMenu(ctx, id)

	return
}

// Update implements intf.FoodUsecase.
func (u *FoodUsecase) Update(c context.Context, request *entity.Food, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.FoodRepository.Update(ctx, request, id)

	return
}

// Delete implements intf.FoodUsecase.
func (u *FoodUsecase) Delete(c context.Context, id string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.ContextTimeout)
	defer cancel()

	err = u.FoodRepository.Delete(ctx, id)

	return
}
