package main

import (
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/middleware"
	"awesomeProject1/internal/store"
	"log"
	"net/http"
)

func main() {
	s := store.New()
	h := handlers.New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.Tasks)

	var handler http.Handler = mux
	handler = middleware.APIKey("secret12345")(handler)
	handler = middleware.Logging("{{message}}")(handler) // можешь заменить на свой текст

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("server started on :8080")
	log.Fatal(server.ListenAndServe())
}
