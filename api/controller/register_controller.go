package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterController struct {
	RegisterUsecase intf.RegisterUsecase
	Env             *bootstrap.Env
}

func (rc *RegisterController) Register(c *gin.Context) {
	var request model.RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	if _, err := rc.RegisterUsecase.FindByEmail(c, request.Email); err == nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "User already exists with the given email"})
		return
	}

	if _, err := rc.RegisterUsecase.FindByPhone(c, request.Phone); err == nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "User already exists with the given phone"})
		return
	}

	hashed, err := helpers.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	current := helpers.Now()
	objectID := primitive.NewObjectID()

	user := entity.User{
		ID:        objectID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  hashed,
		Phone:     request.Phone,
		CreatedAt: *current,
		UpdatedAt: *current,
		UserID:    objectID.Hex(),
	}

	err = rc.RegisterUsecase.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	accessToken, err := rc.RegisterUsecase.CreateAccessToken(&user, rc.Env.AccessTokenSecret, rc.Env.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	refreshToken, err := rc.RegisterUsecase.CreateRefreshToken(&user, rc.Env.RefreshTokenSecret, rc.Env.RefreshTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	response := model.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}
