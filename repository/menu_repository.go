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

type MenuRepository struct {
	Database   mongodb.Database
	Collection string
}

func NewMenuRepository(db mongodb.Database, coll string) intf.MenuRepository {
	return &MenuRepository{
		Database:   db,
		Collection: coll,
	}
}

func (r *MenuRepository) Create(c context.Context, request *entity.Menu) (err error) {
	collection := r.Database.Collection(r.Collection)

	request.ID = primitive.NewObjectID()

	current := helpers.Now()
	request.CreatedAt = *current
	request.UpdatedAt = *current

	_, err = collection.InsertOne(c, request)

	return
}

func (r *MenuRepository) Find(c context.Context) (menus []entity.Menu, err error) {
	// select collection
	collection := r.Database.Collection(r.Collection)

	// get all data
	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		return
	}

	// extract data into variable
	if err = cursor.All(c, &menus); err != nil {
		return
	}

	return
}

func (r *MenuRepository) FindOne(c context.Context, id string) (menu entity.Menu, err error) {
	collection := r.Database.Collection(r.Collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	if err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&menu); err != nil {
		return
	}

	return
}

func (r *MenuRepository) Update(c context.Context, request *entity.Menu, id string) (err error) {
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

func (r *MenuRepository) Delete(c context.Context, id string) (err error) {
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
