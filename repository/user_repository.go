package repository

import (
	"context"

	"github.com/muhafs/go-restaurant-management/domain/entity"
	"github.com/muhafs/go-restaurant-management/domain/intf"
	"github.com/muhafs/go-restaurant-management/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	Database   mongodb.Database
	Collection string
}

func NewUserRepository(db mongodb.Database, coll string) intf.UserRepository {
	return &UserRepository{
		Database:   db,
		Collection: coll,
	}
}

func (ur *UserRepository) Create(c context.Context, user *entity.User) (err error) {
	collection := ur.Database.Collection(ur.Collection)

	_, err = collection.InsertOne(c, user)

	return
}

func (ur *UserRepository) Find(c context.Context) (users []entity.User, err error) {
	collection := ur.Database.Collection(ur.Collection)

	opts := options.Find().SetProjection(bson.D{{Key: "password", Value: 0}})

	cursor, err := collection.Find(c, bson.D{}, opts)
	if err != nil {
		return
	}

	if err = cursor.All(c, &users); err != nil {
		return
	}

	return
}

func (ur *UserRepository) FindOne(c context.Context, id string) (user entity.User, err error) {
	collection := ur.Database.Collection(ur.Collection)

	idHex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	err = collection.FindOne(c, bson.M{"_id": idHex}).Decode(&user)

	return
}

func (ur *UserRepository) FindByEmail(c context.Context, email string) (user entity.User, err error) {
	collection := ur.Database.Collection(ur.Collection)

	err = collection.FindOne(c, bson.M{"email": email}).Decode(&user)

	return
}
