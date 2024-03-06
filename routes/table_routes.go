package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func TableRoutes(r *gin.Engine) {
	r.POST("/tables", controllers.CreateTable)
	r.GET("/tables", controllers.FindTables)
	r.GET("/tables/:table_id", controllers.FindTable)
	r.PATCH("/tables/:table_id", controllers.UpdateTable)
}
