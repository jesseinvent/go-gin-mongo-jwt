package database

import (
	"fmt"
	"time"
	"context"
	"log"
	_"os"

	"go-gin-mongo-jwt/configs"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect() *mongo.Client {
	
	config := configs.Load();

	db_url := config.MONGODB_URL;

	client, err := mongo.NewClient(options.Client().ApplyURI(db_url)); 

	if err != nil {
		log.Fatal(err);
	}

	ctx, cancle := context.WithTimeout(context.Background(), 10 * time.Second);
	 
	defer cancle();

	err = client.Connect(ctx);

	if err != nil {
		log.Fatal(err);
	}

	err = client.Ping(ctx, readpref.Primary());

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB ...");

	return client;
}

var Client *mongo.Client = Connect();

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("go-auth").Collection(collectionName);

	return collection;
}