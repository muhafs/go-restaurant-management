package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionOrder = "orders"

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	OrderDate time.Time          `json:"order_date" bson:"order_date"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	OrderID   string             `json:"order_id" bson:"order_id"`
	TableID   string             `json:"table_id" bson:"table_id"`
}
