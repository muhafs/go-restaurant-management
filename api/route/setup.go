package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/api/middleware"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/mongodb"
)

func Setup(env *bootstrap.Env, db mongodb.Database, timeout time.Duration, gin *gin.Engine) {
	// public APIs
	publicRouter := gin.Group("")

	NewRegisterRouter(env, db, timeout, publicRouter)
	NewLoginRouter(env, db, timeout, publicRouter)
	NewRefreshTokenRouter(env, db, timeout, publicRouter)

	// private APIs
	privateRouter := gin.Group("")
	privateRouter.Use(middleware.Authentication(env.AccessTokenSecret))

	NewProfileRouter(env, db, timeout, privateRouter)
	NewMenuRouter(env, db, timeout, privateRouter)
	// NewTaskRouter(env, timeout, db, privateRouter)
}
