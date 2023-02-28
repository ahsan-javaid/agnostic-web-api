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
	collection string
	param      []string
	w          http.ResponseWriter
	r          *http.Request
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

func (c Context) sendHttp400(msg string) {
	c.w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})

	resp["msg"] = msg

	out, err := json.Marshal(resp)
	utils.Check(err)

	c.w.WriteHeader(http.StatusBadRequest)
	c.w.Write(out)
}

func Router(w http.ResponseWriter, r *http.Request) {
	collection := utils.GetCollectionName(r.URL.Path)
	params := utils.GetParams(r.URL.Path)

	ctx := Context{
		collection: collection,
		param:      params,
		w:          w,
		r:          r,
	}

	if ctx.collection == "" {
		ctx.sendHttp200("ok")
		return
	}

  fmt.Println("path:", r.URL.Path)

	if r.URL.Path == "/users/login" {
		handleLogin(ctx)
		return
	}

	switch r.Method {
	case "GET":
		switch len(ctx.param) {
		case 3:
			handleGetById(ctx)
		default:
			handleGet(ctx)
		}
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

func handleGetById(ctx Context) {
	record := bson.M{}
	id, _ := primitive.ObjectIDFromHex(ctx.param[2])
	
	err := db.DB.Collection(ctx.collection).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&record)
	utils.Check(err)

	ctx.sendHttp200(record)
}

func handleGet(ctx Context) {
	filter := bson.M{}
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

func handleLogin(ctx Context) {
	// Decode body 
	var payload map[string]any
	err := json.NewDecoder(ctx.r.Body).Decode(&payload)
	utils.Check(err)

	// Check if email and password
	email := payload["email"]
	password := payload["password"]

	if email == nil || password == nil {
		ctx.sendHttp400("Email and password required")
		return
	}
	
	record := bson.M{}
	err = db.DB.Collection(ctx.collection).FindOne(context.TODO(), bson.M{"email": email }).Decode(&record)
	utils.Check(err)

	var hash string = fmt.Sprint(record["password"])

	match := utils.CheckPasswordHash(fmt.Sprint(password), hash)

	if match {
		ctx.sendHttp200(record)
	} else {
		ctx.sendHttp400("Invalid email or password")
	}

}
