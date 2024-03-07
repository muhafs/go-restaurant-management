package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

type RefreshTokenController struct {
	RefreshTokenUsecase intf.RefreshTokenUsecase
	Env                 *bootstrap.Env
}

func (rc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var request model.RefreshTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	id, err := rc.RefreshTokenUsecase.ExtractIDFromToken(request.RefreshToken, rc.Env.RefreshTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: "User not found"})
		return
	}

	user, err := rc.RefreshTokenUsecase.FindOne(c, id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Message: "User not found"})
		return
	}

	accessToken, err := rc.RefreshTokenUsecase.CreateAccessToken(&user, rc.Env.AccessTokenSecret, rc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	refreshToken, err := rc.RefreshTokenUsecase.CreateRefreshToken(&user, rc.Env.RefreshTokenSecret, rc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	response := model.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}
