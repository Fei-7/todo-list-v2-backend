package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var cancelCtx context.CancelFunc
var client *mongo.Client
var err error
var UsersCollection *mongo.Collection

func Connect() {
	ctx, cancelCtx = context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		panic(err)
	}

	database := client.Database("todo")
	UsersCollection = database.Collection("users")

	fmt.Println("Database connected")
}

func Disconnect() {
	cancelCtx()
	client.Disconnect(ctx)
	fmt.Println("Database disconnected")
}
