package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func ItemRoutes(r *gin.Engine) {
	r.POST("/items", controllers.CreateItem)
	r.GET("/items", controllers.FindItems)
	r.GET("/items/:item_id", controllers.FindItem)
	r.GET("/items-order/:order_id", controllers.FindItemsByOrder)
	r.PATCH("/items/:item_id", controllers.UpdateItem)
}
