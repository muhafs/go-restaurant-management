package cmd

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
