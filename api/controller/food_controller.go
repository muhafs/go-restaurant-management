package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

type FoodController struct {
	FoodUsecase intf.FoodUsecase
	Env         *bootstrap.Env
}

func (ctl *FoodController) Create(c *gin.Context) {
	var request entity.Food
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	if _, err := ctl.FoodUsecase.FindMenu(c, request.Menu); err == nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "menu not found"})
		return
	}

	if err := ctl.FoodUsecase.Create(c, &request); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{Success: true, Message: "food has created successfully"})
}

func (ctl *FoodController) Find(c *gin.Context) {
	foods, err := ctl.FoodUsecase.Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "food list has fetched successfully", Data: foods})
}

func (ctl *FoodController) FindOne(c *gin.Context) {
	id := c.Param("food_id")

	food, err := ctl.FoodUsecase.FindOne(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "food not found"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "food has fetched successfully", Data: food})
}

func (ctl *FoodController) Update(c *gin.Context) {
	var request entity.Food
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	id := c.Param("food_id")
	if _, err := ctl.FoodUsecase.FindOne(c, id); err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "food not found"})
		return
	}

	if err := ctl.FoodUsecase.Update(c, &request, id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "food has updated successfully"})
}

func (ctl *FoodController) Delete(c *gin.Context) {
	id := c.Param("food_id")
	if _, err := ctl.FoodUsecase.FindOne(c, id); err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "food not found"})
		return
	}

	if err := ctl.FoodUsecase.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "food has deleted successfully"})
}
