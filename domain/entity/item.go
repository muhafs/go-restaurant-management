package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionItem = "items"

type Item struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Quantity  string             `json:"quantity" bson:"quantity"`
	UnitPrice int                `json:"unit_price" bson:"unit_price"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	FoodID    string             `json:"food_id" bson:"food_id"`
	ItemID    string             `json:"item_id" bson:"item_id"`
	OrderID   string             `json:"order_id" bson:"order_id"`
}
