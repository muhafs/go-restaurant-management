package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID        primitive.ObjectID `bson:"_id"`
	Quantity  *string            `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice *int               `json:"unit_price" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	FoodID    *string            `json:"food_id" validate:"required"`
	ItemID    string             `json:"item_id"`
	OrderID   string             `json:"order_id" validate:"required"`
}

type ItemPack struct {
	TableID *string `json:"table_id"`
	Items   []Item  `json:"items"`
}
