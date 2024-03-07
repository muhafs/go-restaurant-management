package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionTable = "tables"

type Table struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	NumberOfGuests int                `json:"number_of_guests" bson:"number_of_guests"`
	TableNumber    int                `json:"table_number" bson:"table_number"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	TableID        string             `json:"table_id" bson:"table_id"`
}
