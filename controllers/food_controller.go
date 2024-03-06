package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/helpers"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")
var validate = validator.New()

func FindFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * limit
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "_id", Value: "null"},
		}},
		{Key: "total_count", Value: bson.D{
			{Key: "$sum", Value: 1},
		}},
		{Key: "data", Value: bson.D{
			{Key: "$push", Value: "$$ROOT"},
		}},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "food_items", Value: bson.D{
			{Key: "$slice", Value: []interface{}{
				"$data", startIndex, limit,
			}},
		}},
	}}}

	result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var foods []bson.M
	if err := result.All(ctx, &foods); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foods[0])
}

func FindFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	foodID := c.Param("food_id")
	var food models.Food

	if err := foodCollection.FindOne(ctx, bson.M{"foodid": foodID}).Decode(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error occured while fetching the food item",
		})
		return
	}

	c.JSON(http.StatusOK, food)
}

func CreateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var menu models.Menu
	var food models.Food

	// catch request body
	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// validate
	if vErr := validate.Struct(food); vErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
		return
	}

	// check for menu
	if err := menuCollection.FindOne(ctx, bson.M{"menuid": food.MenuID}).Decode(&menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "menu not found"})
		return
	}

	// set creation time
	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// set new ID
	food.ID = primitive.NewObjectID()
	food.FoodID = food.ID.Hex()

	// set price precision
	var p = helpers.ToFixed(*food.Price, 2)
	food.Price = &p

	// create food record
	result, err := foodCollection.InsertOne(ctx, food)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "food item failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	foodID := c.Param("food_id")

	var menu models.Menu
	var food models.Food

	if err := c.BindJSON(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if food.Name != nil {
		updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
	}

	if food.Price != nil {
		updateObj = append(updateObj, bson.E{Key: "price", Value: food.Price})
	}

	if food.FoodImage != nil {
		updateObj = append(updateObj, bson.E{Key: "food_image", Value: food.FoodImage})
	}

	if food.MenuID != nil {
		if err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuID}).Decode(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "menu not found"})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "menu_id", Value: food.MenuID})
	}

	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

	result, err := foodCollection.UpdateOne(
		ctx,
		bson.M{"food_id": foodID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "food item failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
}
