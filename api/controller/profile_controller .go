package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

type ProfileController struct {
	ProfileUsecase intf.ProfileUsecase
	Env            *bootstrap.Env
}

func (pc *ProfileController) Find(c *gin.Context) {
	profiles, err := pc.ProfileUsecase.Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profiles)
}

func (pc *ProfileController) FindOne(c *gin.Context) {
	userID := c.GetString("x-user-id")

	profile, err := pc.ProfileUsecase.FindOne(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
