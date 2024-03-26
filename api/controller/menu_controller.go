package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/bootstrap"
	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/domain/model"
)

type MenuController struct {
	MenuUsecase intf.MenuUsecase
	Env         *bootstrap.Env
}

func (ctl *MenuController) Create(c *gin.Context) {
	var request entity.Menu
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	if err := ctl.MenuUsecase.Create(c, &request); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{Success: true, Message: "menu has created successfully"})
}

func (ctl *MenuController) Find(c *gin.Context) {
	menus, err := ctl.MenuUsecase.Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "menu list has fetched successfully", Data: menus})
}

func (ctl *MenuController) FindOne(c *gin.Context) {
	id := c.Param("menu_id")

	menu, err := ctl.MenuUsecase.FindOne(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "menu not found"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "menu has fetched successfully", Data: menu})
}

func (ctl *MenuController) Update(c *gin.Context) {
	var request entity.Menu
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	id := c.Param("menu_id")
	if _, err := ctl.MenuUsecase.FindOne(c, id); err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "menu not found"})
		return
	}

	if err := ctl.MenuUsecase.Update(c, &request, id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "menu has updated successfully"})
}

func (ctl *MenuController) Delete(c *gin.Context) {
	id := c.Param("menu_id")
	if _, err := ctl.MenuUsecase.FindOne(c, id); err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Success: false, Message: "menu not found"})
		return
	}

	if err := ctl.MenuUsecase.Delete(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Success: true, Message: "menu has deleted successfully"})
}
