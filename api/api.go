package api

import (
	"context"
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	db "agnostic-web-api/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("method", r.Method)
	switch r.Method {
	case "GET":
		handleGet(w, r)
	case "PUT":
		handlePut(w, r)
	case "POST":
		handlePost(w, r)
	case "DELETE":
		handleGet(w, r)
	default:
		fmt.Fprintf(w, "METHOD NOT SUPPORTED")
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path is", r.URL.Path)
	urlParts := strings.Split(r.URL.Path, "/")

	collectionName := urlParts[1]

	fmt.Println("c name", collectionName)

	dbClient := db.DB
	coll := dbClient.Collection(collectionName)
	title := "Back to the Future"

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", jsonData)

  w.Write(jsonData)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		panic(err)
	}
	dbClient := db.DB
	result, err := dbClient.Collection("users").InsertOne(context.TODO(), payload)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)

}

func handlePut(w http.ResponseWriter, r *http.Request) {

}

func handleDelete(w http.ResponseWriter, r *http.Request) {

}
