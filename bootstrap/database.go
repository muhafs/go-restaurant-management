package bootstrap

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muhafs/go-restaurant-management/mongodb"
)

func NewMongoDatabase(env *Env) mongodb.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPass

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongodb.NewClient(ctx, mongodbURI)
	if err != nil {
		log.Fatal(err)
		log.Fatal("error on connection")
	}

	err = client.Ping(ctx)
	if err != nil {
		log.Fatal(err)
		log.Fatal("error on ping")
	}

	return client
}

func CloseMongoDBConnection(client mongodb.Client) {
	if client == nil {
		return
	}

	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection to MongoDB closed.")
}
