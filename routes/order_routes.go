package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func OrderRoutes(r *gin.Engine) {
	r.POST("/orders", controllers.CreateOrder)
	r.GET("/orders", controllers.FindOrders)
	r.GET("/orders/:order_id", controllers.FindOrder)
	r.PATCH("/orders/:order_id", controllers.UpdateOrder)
}
