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

func NewRegisterRouter(env *bootstrap.Env, db mongodb.Database, timeout time.Duration, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, entity.CollectionUser)
	ru := usecase.NewRegisterUsecase(ur, timeout)
	rc := controller.RegisterController{
		RegisterUsecase: ru,
		Env:             env,
	}

	group.POST("/register", rc.Register)
}
