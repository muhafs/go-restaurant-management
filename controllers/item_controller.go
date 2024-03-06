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

var itemCollection *mongo.Collection = database.OpenCollection(database.Client, "item")

func FindItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	result, err := itemCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var items []bson.M
	if err := result.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func FindItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	itemID := c.Param("item_id")
	var item models.Item

	if err := itemCollection.FindOne(ctx, bson.M{"itemid": itemID}).Decode(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func FindItemsByOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	orderID := c.Param("order_id")

	items, err := ItemsByOrder(ctx, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func ItemsByOrder(c context.Context, orderID string) (items []bson.M, err error) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "order_id", Value: orderID},
	}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: "food_id"},
		{Key: "as", Value: "food"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$food"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	orderLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "order"},
		{Key: "localField", Value: "order_id"},
		{Key: "foreignField", Value: "order_id"},
		{Key: "as", Value: "order"},
	}}}
	orderUnwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$order"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	tableLookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "table"},
		{Key: "localField", Value: "order.table_id"},
		{Key: "foreignField", Value: "table_id"},
		{Key: "as", Value: "table"},
	}}}
	tableUnwindStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$table"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	projectStage := bson.D{{Key: "$porject", Value: bson.D{
		{Key: "id", Value: 0},
		{Key: "amount", Value: "$food.price"},
		{Key: "total_count", Value: 1},
		{Key: "food_name", Value: "$food.name"},
		{Key: "food_image", Value: "$food.food_image"},
		{Key: "table_number", Value: "$table.table_number"},
		{Key: "table_id", Value: "$table.table_id"},
		{Key: "order_id", Value: "$order.order_id"},
		{Key: "price", Value: "$food.price"},
		{Key: "quantity", Value: 1},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "order_id", Value: "$order.order_id"},
			{Key: "table_id", Value: "$table.table_id"},
			{Key: "table_number", Value: "$table.table_number"},
		}},
		{Key: "payment_due", Value: bson.D{
			{Key: "$sum", Value: "$amount"},
		}},
		{Key: "total_count", Value: bson.D{
			{Key: "$sum", Value: 1},
		}},
		{Key: "items", Value: bson.D{
			{Key: "$push", Value: "$$ROOT"},
		}},
	}}}

	lastStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "id", Value: 0},
		{Key: "payment_due", Value: 1},
		{Key: "total_count", Value: 1},
		{Key: "table_number", Value: "$_id.table_number"},
		{Key: "items", Value: 1},
	}}}

	result, err := itemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		orderLookupStage,
		orderUnwindStage,
		tableLookupStage,
		tableUnwindStage,
		projectStage,
		groupStage,
		lastStage,
	})
	if err != nil {
		return
	}

	if err = result.All(ctx, &items); err != nil {
		return
	}

	return
}

func CreateItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	var itemPack models.ItemPack
	if err := c.BindJSON(&itemPack); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	order.TableID = itemPack.TableID
	order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderID := ItemCreator(ctx, order)

	var newItems []any
	for _, item := range itemPack.Items {
		item.OrderID = orderID

		if vErr := validate.Struct(item); vErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": vErr.Error()})
			return
		}

		item.ID = primitive.NewObjectID()
		item.ItemID = item.ID.Hex()

		item.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		item.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		p := helpers.Round(helpers.ToFixed(float64(*item.UnitPrice), 2))
		item.UnitPrice = &p

		newItems = append(newItems, item)
	}

	result, err := itemCollection.InsertMany(ctx, newItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "items failed to create"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	itemID := c.Param("item_id")

	var item models.Item
	if err := c.BindJSON(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var updateObj bson.D

	if item.UnitPrice != nil {
		updateObj = append(updateObj, bson.E{Key: "unit_price", Value: item.UnitPrice})
	}

	if item.Quantity != nil {
		updateObj = append(updateObj, bson.E{Key: "quantity", Value: item.Quantity})
	}

	if item.FoodID != nil {
		updateObj = append(updateObj, bson.E{Key: "food_id", Value: item.FoodID})
	}

	item.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: item.UpdatedAt})

	result, err := itemCollection.UpdateOne(
		ctx,
		bson.M{"item_id": itemID},
		bson.D{
			{Key: "$set", Value: updateObj},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "item failed to update"})
		return
	}

	c.JSON(http.StatusOK, result)
}
