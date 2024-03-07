package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionFood = "foods"

type Food struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Price     float64            `json:"price" bson:"price"`
	FoodImage string             `json:"food_image" bson:"food_image"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	FoodID    string             `json:"food_id" bson:"food_id"`
	MenuID    string             `json:"menu_id" bson:"menu_id"`
}
