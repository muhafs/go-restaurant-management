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

func NewLoginRouter(env *bootstrap.Env, db mongodb.Database, timeout time.Duration, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, entity.CollectionUser)
	lu := usecase.NewLoginUsecase(ur, timeout)
	lc := controller.LoginController{
		LoginUsecase: lu,
		Env:          env,
	}

	group.POST("/login", lc.Login)
}
