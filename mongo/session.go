package mongo

import "go.mongodb.org/mongo-driver/mongo"

type mongoSession struct {
	mongo.Session
}
