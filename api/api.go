package api

import (
	db "agnostic-web-api/db"
	utils "agnostic-web-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func Router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

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
	collection := utils.GetCollectionName(r.URL.Path)
	cursor, err := db.DB.Collection(collection).Find(context.TODO(), bson.M{})
	utils.Check(err)
	
	var results []bson.M
	
	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.Check(err)
	}

	defer cursor.Close(context.TODO())
	
	if err := cursor.Err(); err != nil {
		utils.Check(err)
	}

	out, err := json.Marshal(results)
	utils.Check(err)
	
	w.Write(out)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any
	err := json.NewDecoder(r.Body).Decode(&payload)
	utils.Check(err)

	collection := utils.GetCollectionName(r.URL.Path)
	result, err := db.DB.Collection(collection).InsertOne(context.TODO(), payload)

	utils.Check(err)

	out, err := json.Marshal(result)
	utils.Check(err)

  w.Write(out)
}

func handlePut(w http.ResponseWriter, r *http.Request) {

}

func handleDelete(w http.ResponseWriter, r *http.Request) {

}
