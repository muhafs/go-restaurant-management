package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/helpers"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func FindMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	result, err := menuCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var menus []bson.M
	if err := result.All(ctx, &menus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menus)
}

func FindMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	menuID := c.Param("menu_id")

	var menu models.Menu
	if err := menuCollection.FindOne(ctx, bson.M{"menuid": menuID}).Decode(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func CreateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if vErr := validate.Struct(menu); vErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": vErr.Error()})
		return
	}

	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	menu.ID = primitive.NewObjectID()
	menu.MenuID = menu.ID.Hex()

	result, err := menuCollection.InsertOne(ctx, menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "menu item failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateMenu(c *gin.Context) {
	// create timeout
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	// extract user request into a model
	var menu models.Menu
	if err := c.BindJSON(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// extract id from URL
	menuID := c.Param("menu_id")

	if menu.StartDate == nil && menu.EndDate == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Start and End date are empty"})
		return
	}

	if !helpers.InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kindly retype the time"})
		return
	}

	var updateObj bson.D
	updateObj = append(updateObj, bson.E{Key: "start_date", Value: menu.StartDate})
	updateObj = append(updateObj, bson.E{Key: "end_date", Value: menu.EndDate})

	if menu.Name != "" {
		updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
	}

	if menu.Category != "" {
		updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
	}

	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

	result, err := menuCollection.UpdateOne(
		ctx,
		bson.M{"menu_id": menuID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "menu item failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
}
