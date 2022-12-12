package api

import (
	db "agnostic-web-api/db"
	utils "agnostic-web-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Context struct {
	collection   string
	param        []string
	w            http.ResponseWriter
	r            *http.Request
}

func (c Context) send(data []byte) {
	c.w.Header().Set("Content-Type", "application/json")
	c.w.Write(data)
}


func Router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	collection := utils.GetCollectionName(r.URL.Path)
	params := utils.GetParams(r.URL.Path)

	c := Context {
		collection:   collection,
		param:    params,
		w: w,
		r: r,
  }

	switch r.Method {
	case "GET":
		handleGet(c)
	case "PUT":
		handlePut(w, r)
	case "POST":
		handlePost(w, r)
	case "DELETE":
		handleDelete(w, r)
	default:
		fmt.Fprintf(w, "METHOD NOT SUPPORTED")
	}
}

func handleGet(c Context) {
	cursor, err := db.DB.Collection(c.collection).Find(context.TODO(), bson.M{})
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

	c.send(out)
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
	var payload map[string]any
	err := json.NewDecoder(r.Body).Decode(&payload)
	utils.Check(err)

	collection := utils.GetCollectionName(r.URL.Path)
	id := utils.GetURLParam(r.URL.Path)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": payload,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	result := db.DB.Collection(collection).FindOneAndUpdate(context.TODO(), filter, update, &opt)
	utils.Check(result.Err())
	doc := bson.M{}
	decodeErr := result.Decode(&doc)
	utils.Check(decodeErr)
	out, err := json.Marshal(doc)
	utils.Check(err)
	w.Write(out)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	collection := utils.GetCollectionName(r.URL.Path)
	id := utils.GetURLParam(r.URL.Path)

	where := bson.M{"_id": id}
	res, err := db.DB.Collection(collection).DeleteOne(context.TODO(), where)
	utils.Check(err)

	out, err := json.Marshal(res)
	utils.Check(err)
	w.Write(out)
}
