package mongo

import "go.mongodb.org/mongo-driver/mongo"

type SingleResult interface {
	Decode(interface{}) error
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}
