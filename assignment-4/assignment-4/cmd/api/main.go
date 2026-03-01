package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Genre      string `json:"genre"`
	Budget     int    `json:"budget"`
	Hero       string `json:"hero"`
	Heroine    string `json:"heroine"`
}

var db *sql.DB

func connectDB() *sql.DB {
	host := getEnv("DB_HOST", "db")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "moviesdb")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var database *sql.DB
	var err error

	// Retry connection — wait for DB healthcheck
	for i := 0; i < 10; i++ {
		database, err = sql.Open("postgres", dsn)
		if err == nil {
			err = database.Ping()
		}
		if err == nil {
			log.Println("Successfully connected to the database!")
			return database
		}
		log.Printf("Waiting for database... attempt %d/10: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("Could not connect to database: %v", err)
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// GET /movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, genre, budget, hero, heroine FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		rows.Scan(&m.ID, &m.Title, &m.Genre, &m.Budget, &m.Hero, &m.Heroine)
		movies = append(movies, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id}
func getMovie(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var m Movie
	err := db.QueryRow("SELECT id, title, genre, budget, hero, heroine FROM movies WHERE id=$1", id).
		Scan(&m.ID, &m.Title, &m.Genre, &m.Budget, &m.Hero, &m.Heroine)
	if err == sql.ErrNoRows {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// POST /movies
func createMovie(w http.ResponseWriter, r *http.Request) {
	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := db.QueryRow(
		"INSERT INTO movies (title, genre, budget, hero, heroine) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		m.Title, m.Genre, m.Budget, m.Hero, m.Heroine,
	).Scan(&m.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(m)
}

// PUT /movies/{id}
func updateMovie(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var m Movie
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	result, err := db.Exec(
		"UPDATE movies SET title=$1, genre=$2, budget=$3, hero=$4, heroine=$5 WHERE id=$6",
		m.Title, m.Genre, m.Budget, m.Hero, m.Heroine, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	m.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// DELETE /movies/{id}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	result, err := db.Exec("DELETE FROM movies WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
}

func main() {
	db = connectDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// Graceful shutdown (EASY optional feature)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down gracefully...")
		db.Close()
		os.Exit(0)
	}()

	fmt.Println("Starting the Server...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
