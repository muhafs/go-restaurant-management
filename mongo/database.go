package mongo

import "go.mongodb.org/mongo-driver/mongo"

type Database interface {
	Collection(string) Collection
	Client() Client
}

type mongoDatabase struct {
	db *mongo.Database
}

func (md *mongoDatabase) Collection(colName string) Collection {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.db.Client()
	return &mongoClient{cl: client}
}
