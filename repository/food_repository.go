package repository

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/internal/helpers"
	"github.com/muhafs/go-restaurant-management/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FoodRepository struct {
	Database   mongodb.Database
	Collection string
	Foreign    string
}

func NewFoodRepository(db mongodb.Database, coll string, foreign string) intf.FoodRepository {
	return &FoodRepository{
		Database:   db,
		Collection: coll,
		Foreign:    foreign,
	}
}

func (r *FoodRepository) Create(c context.Context, request *entity.Food) (err error) {
	collection := r.Database.Collection(r.Collection)

	request.ID = primitive.NewObjectID()

	current := helpers.Now()
	request.CreatedAt = *current
	request.UpdatedAt = *current

	_, err = collection.InsertOne(c, request)

	return
}

func (r *FoodRepository) Find(c context.Context) (foods []entity.Food, err error) {
	// select collection
	collection := r.Database.Collection(r.Collection)

	// get all data
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return
	}

	// extract data into variable
	if err = cursor.All(c, &foods); err != nil {
		return
	}

	return
}

func (r *FoodRepository) FindOne(c context.Context, id string) (food entity.Food, err error) {
	collection := r.Database.Collection(r.Collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	if err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&food); err != nil {
		return
	}

	return
}

func (r *FoodRepository) Update(c context.Context, request *entity.Food, id string) (err error) {
	collection := r.Database.Collection(r.Collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	filter := bson.M{"_id": idHex}
	opt := options.Update().SetUpsert(true)

	request.UpdatedAt = *helpers.Now()
	update := bson.D{{Key: "$set", Value: request}}

	if _, err = collection.UpdateOne(c, filter, update, opt); err != nil {
		return
	}

	return
}

func (r *FoodRepository) FindMenu(c context.Context, id string) (menu entity.Menu, err error) {
	foreign := r.Database.Collection(r.Foreign)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	if err = foreign.FindOne(c, bson.M{"_id": idHex}).Decode(&menu); err != nil {
		return
	}

	return
}

func (r *FoodRepository) Delete(c context.Context, id string) (err error) {
	collection := r.Database.Collection(r.Collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	filter := bson.M{"_id": idHex}

	if _, err = collection.DeleteOne(c, filter); err != nil {
		return
	}

	return
}
