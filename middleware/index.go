package middleware

import (
	"net/http"
	"os"
	"fmt"
)

// Todo: middleware
func Index(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Listening on port: %s", os.Getenv("PORT"))))
		h.ServeHTTP(w, r)
	})
}
