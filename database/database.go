package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context
var cancelCtx context.CancelFunc
var client *mongo.Client
var err error
var mongoURI string
var UsersCollection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	mongoURI = os.Getenv("MONGODB_URI")
}

func Connect() {
	ctx, cancelCtx = context.WithTimeout(context.Background(), 10*time.Second)
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
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
