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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func FindOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	result, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var orders []bson.M
	if err := result.All(ctx, &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
	return
}

func FindOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	orderID := c.Param("order_id")
	var order models.Order

	if err := orderCollection.FindOne(ctx, bson.M{"orderid": orderID}).Decode(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
	return
}

func CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var table models.Table
	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if vErr := validate.Struct(order); vErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
		return
	}

	if err := tableCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "table not found"})
		return
	}

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	result, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
	return
}

func UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	orderID := c.Param("order_id")

	var table models.Table
	var order models.Order

	if err := c.BindJSON(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if order.TableID != nil {
		if err := orderCollection.FindOne(ctx, bson.M{"table_id": order.TableID}).Decode(&table); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "table not found"})
			return
		}

		updateObj = append(updateObj, bson.E{Key: "table_id", Value: order.TableID})
	}

	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

	result, err := orderCollection.UpdateOne(
		ctx,
		bson.M{"order_id": orderID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "order failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
	return
}

func ItemCreator(c context.Context, order models.Order) string {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	orderCollection.InsertOne(ctx, order)
	return order.OrderID
}
