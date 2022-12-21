package api

import (
	db "agnostic-web-api/db"
	utils "agnostic-web-api/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Context struct {
	collection   string
	param        []string
	w            http.ResponseWriter
	r            *http.Request
}

func (c Context) sendHttp200(data interface{}) {
	c.w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})

	resp["data"] = data

	out, err := json.Marshal(resp)
	utils.Check(err)

	c.w.WriteHeader(http.StatusOK)
	c.w.Write(out)
}


func Router(w http.ResponseWriter, r *http.Request) {
	collection := utils.GetCollectionName(r.URL.Path)
	params := utils.GetParams(r.URL.Path)

	ctx := Context {
		collection:   collection,
		param:    params,
		w: w,
		r: r,
  }

	switch r.Method {
	case "GET":
		handleGet(ctx)
	case "PUT":
		handlePut(ctx)
	case "POST":
		handlePost(ctx)
	case "DELETE":
		handleDelete(ctx)
	default:
		fmt.Fprintf(w, "METHOD NOT SUPPORTED")
	}
}

func handleGet(ctx Context) {
	filter := bson.M{}
	if len(ctx.param) == 3 {
		record := bson.M{}
		filter = bson.M{"_id": ctx.param[2]}
		id, _ := primitive.ObjectIDFromHex(ctx.param[2])
		err := db.DB.Collection(ctx.collection).FindOne(context.TODO(),bson.M{"_id": id}).Decode(&record)
		utils.Check(err)
		ctx.sendHttp200(record)
		return
	}

	cursor, err := db.DB.Collection(ctx.collection).Find(context.TODO(), filter)
	utils.Check(err)

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		utils.Check(err)
	}

	defer cursor.Close(context.TODO())

	if err := cursor.Err(); err != nil {
		utils.Check(err)
	}

	ctx.sendHttp200(results)
}

func handlePost(ctx Context) {
	var payload map[string]any
	err := json.NewDecoder(ctx.r.Body).Decode(&payload)
	utils.Check(err)

	collection := utils.GetCollectionName(ctx.r.URL.Path)
	result, err := db.DB.Collection(collection).InsertOne(context.TODO(), payload)

	utils.Check(err)

	ctx.sendHttp200(result)
}

func handlePut(ctx Context) {
	var payload map[string]any
	err := json.NewDecoder(ctx.r.Body).Decode(&payload)
	utils.Check(err)

	collection := utils.GetCollectionName(ctx.r.URL.Path)
	id := ctx.param[2]

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

	ctx.sendHttp200(doc)
}

func handleDelete(ctx Context) {
	collection := utils.GetCollectionName(ctx.r.URL.Path)
	id := ctx.param[2]

	where := bson.M{"_id": id}
	res, err := db.DB.Collection(collection).DeleteOne(context.TODO(), where)
	utils.Check(err)

	ctx.sendHttp200(res)
}
