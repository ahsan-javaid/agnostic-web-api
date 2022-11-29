

package main

import (
    "log"
    "fmt"
    "os"
    "net/http"
		api "agnostic-web-api/api"
)

func main() {
    http.HandleFunc("/", api.Router)

    PORT := os.Getenv("PORT")

		if PORT == "" {
			PORT = "5000"
		}

    log.Println("Listening on port ", PORT)

    log.Fatal(http.ListenAndServe(fmt.Sprint(":", PORT), nil))
}