package Mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strings"
)

type Client struct {
	client *mongo.Client
}

type Db struct {
	db *mongo.Database
}

func (client Client) GetDB(name string) *Db {
	return &Db{
		db: client.client.Database(name),
	}
}

func (db Db) GetCollection(name string) *Collection {
	return &Collection{
		collection: db.db.Collection(name),
	}
}

func Connect(username string, password string, clusterHost string, ctx context.Context) *Client {
	if username == "" || password == "" || clusterHost == "" {
		log.Fatal("MongoDb Connect: Invalid MongoDB credentials Provided")
	}
	uri := "mongodb+srv://" + username + ":" + strings.Replace(password, " ", "%20", -1) + "@" + clusterHost + "/?retryWrites=true&w=majority"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDb Connect: Unable to connect to MongoDB.", err)
	}

	return &Client{
		client: client,
	}
}

func (client Client) Disconnect(ctx context.Context) {
	if err := client.client.Disconnect(ctx); err != nil {
		log.Println("MongoDb Disconnect: Unable to disconnect from MongoDB.", err)
	} else {
		log.Println("MongoDb Disconnect: Successfully disconnected from MongoDB.")
	}
}
