package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionMenu = "menus"

type Menu struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" binding:"required"`
	Category  string             `json:"category" bson:"category" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
