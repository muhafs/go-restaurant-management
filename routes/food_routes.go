package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/controllers"
)

func FoodRoutes(r *gin.Engine) {
	r.POST("/foods", controllers.CreateFood)
	r.GET("/foods", controllers.FindFoods)
	r.GET("/foods/:food_id", controllers.FindFood)
	r.PATCH("/foods/:food_id", controllers.UpdateFood)
}
