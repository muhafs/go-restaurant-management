package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/domain/model"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
)

func Authentication(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: "No Authorization header provided"})
			c.Abort()
			return
		}

		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: "Malformed token"})
			c.Abort()
		}

		tokenString := tokens[1]
		authorized, _ := helpers.IsAuthorized(tokenString, secret)
		if !authorized {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: "Not authorized"})
			c.Abort()
			return
		}

		userID, err := helpers.ExtractIDFromToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: err.Error()})
			c.Abort()
			return
		}

		c.Set("x-user-id", userID)
		c.Next()
	}

}
