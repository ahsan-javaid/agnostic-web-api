package db

import (
	"context"
	"fmt"
	"os"
	"log"

	utils "agnostic-web-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func Connect() {
	URI := os.Getenv("DB_URI")
	DB_NAME := os.Getenv("DB_NAME")

	connectionUri := fmt.Sprintf("%s/%s", URI, DB_NAME)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionUri))
	utils.Check(err)
	log.Println("Connected to mongodb: ", connectionUri)
  DB = client.Database(DB_NAME)
}
