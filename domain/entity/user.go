package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionUser = "users"

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    string             `json:"first_name" bson:"first_name"`
	LastName     string             `json:"last_name" bson:"last_name"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Phone        string             `json:"phone" bson:"phone"`
	Avatar       string             `json:"avatar" bson:"avatar"`
	Token        string             `json:"token" bson:"token"`
	RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
	UserID       string             `json:"user_id" bson:"user_id"`
}
