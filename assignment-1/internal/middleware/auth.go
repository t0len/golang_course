package middleware

import (
	"encoding/json"
	"net/http"
)

type errResp struct {
	Error string `json:"error"`
}

func APIKey(required string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			key := r.Header.Get("X-API-KEY")
			if key != required {
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(errResp{Error: "unauthorized"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
