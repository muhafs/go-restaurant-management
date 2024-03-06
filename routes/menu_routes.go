package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func MenuRoutes(r *gin.Engine) {
	r.POST("/menus", controllers.CreateMenu)
	r.GET("/menus", controllers.FindMenus)
	r.GET("/menus/:menu_id", controllers.FindMenu)
	r.PATCH("/menus/:menu_id", controllers.UpdateMenu)
}
