package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionFood = "foods"

type Food struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name" binding:"required,min=2,max=100"`
	Price     int64              `json:"price" bson:"price" binding:"required"`
	FoodImage string             `json:"food_image,omitempty" bson:"food_image,omitempty"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Menu      string             `json:"menu_id" bson:"menu_id" binding:"required"`
}
