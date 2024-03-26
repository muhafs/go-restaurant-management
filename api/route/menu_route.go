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

func NewMenuRouter(env *bootstrap.Env, db mongodb.Database, timeout time.Duration, group *gin.RouterGroup) {
	repo := repository.NewMenuRepository(db, entity.CollectionMenu)
	uscs := usecase.NewMenuUsecase(repo, timeout)
	ctl := &controller.MenuController{MenuUsecase: uscs, Env: env}

	group.POST("/menu", ctl.Create)
	group.GET("/menu", ctl.Find)
	group.GET("/menu/:menu_id", ctl.FindOne)
	group.PUT("/menu/:menu_id", ctl.Update)
	group.DELETE("/menu/:menu_id", ctl.Delete)
}
