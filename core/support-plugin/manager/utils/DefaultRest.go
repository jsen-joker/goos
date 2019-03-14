package utils

import (
	"net/http"
)

func HttpPipeRest(inner http.Handler, name string) http.HandlerFunc  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		inner.ServeHTTP(w, r)
	})
}