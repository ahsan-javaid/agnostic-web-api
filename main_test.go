package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"
    api "agnostic-web-api/api"
    db "agnostic-web-api/db"
    config "agnostic-web-api/config"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

func TestAPI(t *testing.T) {
    // Create a request to pass to the API handler
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a response recorder to record the response
    rr := httptest.NewRecorder()

    // Create a handler for the API endpoint
    handler := http.HandlerFunc(api.Router)

    // Call the API endpoint handler with the request and response recorder
    handler.ServeHTTP(rr, req)

    // Check the status code is what we expect
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("API endpoint returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Check the response body is what we expect
    expected := `{"data":"ok"}`
    if rr.Body.String() != expected {
        t.Errorf("API endpoint returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
}

func TestDBConnection(t *testing.T) {
    config.LoadEnv(".env")
   
    db.Connect()
   
    err := db.DB.Client().Ping(context.TODO(), readpref.Primary())

    if err !=nil {
        t.Errorf("Unable to connect to db %v", err)
    }
}