package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"practice5/internal/handler"
	"practice5/internal/repository"

	_ "github.com/lib/pq"
)

func main() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "Aidosmaidos2003123")
	dbname := getEnv("DB_NAME", "practice5")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping db:", err)
	}
	log.Println("Connected to database")

	repo := repository.New(db)
	h := handler.New(repo)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
