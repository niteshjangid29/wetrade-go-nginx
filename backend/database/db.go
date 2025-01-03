package database

import (
	"context"
	"fmt"

	"github.com/wetrade/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(cfg config.Config) *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.MONGO_URI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(cfg.DB_NAME).RunCommand(context.TODO(), bson.D{
		{Key: "ping", Value: 1},
	}).Err(); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return client
}

var Client *mongo.Client = ConnectDB(config.LoadConfig())

func StocksData(collection string, cfg config.Config) *mongo.Collection {
	return Client.Database(cfg.DB_NAME).Collection(collection)
}

func UserData(collection string, cfg config.Config) *mongo.Collection {
	return Client.Database(cfg.DB_NAME).Collection(collection)
}

func ContactData(collection string, cfg config.Config) *mongo.Collection {
	return Client.Database(cfg.DB_NAME).Collection(collection)
}
