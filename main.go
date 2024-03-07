// package main

// import (
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/muhafs/go-restaurant-management/database"
// 	"github.com/muhafs/go-restaurant-management/middleware"
// 	"github.com/muhafs/go-restaurant-management/routes"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// func main() {
// 	port := os.Getenv("PORT")

// 	if port == "" {
// 		port = "8080"
// 	}

// 	router := gin.Default()

// 	routes.UserRoutes(router)

// 	router.Use(middleware.Authentication())

// 	routes.FoodRoutes(router)
// 	routes.MenuRoutes(router)
// 	routes.TableRoutes(router)
// 	routes.OrderRoutes(router)
// 	routes.ItemRoutes(router)
// 	routes.InvoiceRoutes(router)

// 	router.Run(":" + port)
// }

package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/api/route"
	"github.com/muhafs/go-restaurant-management/bootstrap"
)

func main() {
	app := bootstrap.App()

	env := app.Env

	db := app.Mongo.Database(env.DBName)
	defer app.CloseDBConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gr := gin.Default()

	route.Setup(env, db, timeout, gr)

	gr.Run(env.ServerAddress)
}
