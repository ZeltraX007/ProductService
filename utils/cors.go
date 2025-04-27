package utils

import (
	"log"
	"net/http"
)

func CorsFilter(next http.Handler) http.Handler {
	log.Printf("Entered CorsFilter")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		next.ServeHTTP(w, r)
	})
}
