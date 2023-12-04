package main

import (
	api "agnostic-web-api/api"
	// middleware "agnostic-web-api/middleware"
	config "agnostic-web-api/config"
	db "agnostic-web-api/db"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
  config.LoadEnv(".env") // Loading env variables into process from .env file
  db.Connect() // Connecting mongodb

	http.HandleFunc("/", api.Router)
	PORT := os.Getenv("PORT")

	log.Println("Listening on port: ", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprint(":", PORT), nil))
}
