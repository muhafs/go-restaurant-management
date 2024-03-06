package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhafs/go-restaurant-management/database"
	"github.com/muhafs/go-restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func FindTables(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	result, err := tableCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var tables []bson.M
	if err := result.All(ctx, &tables); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

func FindTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	tableID := c.Param("table_id")
	var table models.Table

	if err := tableCollection.FindOne(ctx, bson.M{"tableid": tableID}).Decode(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func CreateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var table models.Table

	if err := c.BindJSON(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if vErr := validate.Struct(table); vErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
		return
	}

	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	table.ID = primitive.NewObjectID()
	table.TableID = table.ID.Hex()

	result, err := tableCollection.InsertOne(ctx, table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "table failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	tableID := c.Param("table_id")

	var table models.Table
	if err := c.BindJSON(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if table.NumberOfGuests != nil {
		updateObj = append(updateObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
	}

	if table.TableNumber != nil {
		updateObj = append(updateObj, bson.E{Key: "table_number", Value: table.TableNumber})
	}

	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

	result, err := tableCollection.UpdateOne(
		ctx,
		bson.M{"table_id": tableID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "table failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
}
