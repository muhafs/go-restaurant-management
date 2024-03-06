package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controllers.FindUsers)
	r.GET("/users/:user_id", controllers.FindUser)

	r.POST("/users/register", controllers.Register)
	r.POST("/users/login", controllers.Login)
}
