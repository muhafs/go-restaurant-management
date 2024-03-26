package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/api/controller"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/mongodb"
	"github.com/muhafs/go-restaurant-management/repository"
	"github.com/muhafs/go-restaurant-management/usecase"
)

func NewFoodRouter(env *bootstrap.Env, db mongodb.Database, timeout time.Duration, group *gin.RouterGroup) {
	repo := repository.NewFoodRepository(db, entity.CollectionFood, entity.CollectionMenu)
	uscs := usecase.NewFoodUsecase(repo, timeout)
	ctl := &controller.FoodController{FoodUsecase: uscs, Env: env}

	group.POST("/food", ctl.Create)
	group.GET("/food", ctl.Find)
	group.GET("/food/:food_id", ctl.FindOne)
	group.PUT("/food/:food_id", ctl.Update)
	group.DELETE("/food/:food_id", ctl.Delete)
}
