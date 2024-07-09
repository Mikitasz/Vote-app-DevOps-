package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var db *sql.DB

type Vote struct {
	Category string `json:"category"`
	Votes    int    `json:"votes"`
}

func main() {
	var err error
	dbUser := os.Getenv("DB_USER")
	fmt.Printf("Database User: %s\n", dbUser)
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = initDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/votes", getVotesHandler)
	http.Handle("/", http.FileServer(http.Dir("template")))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDB() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS vote_app (
		id SERIAL PRIMARY KEY,
		category VARCHAR(50) NOT NULL,
		votes INTEGER NOT NULL DEFAULT 0
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	insertCategoriesQuery := `
	INSERT INTO vote_app (category, votes)
	VALUES
		('dogs', 0),
		('cats', 0)
	ON CONFLICT (category) DO NOTHING;`

	_, err = db.Exec(insertCategoriesQuery)
	if err != nil {
		return fmt.Errorf("failed to insert initial categories: %v", err)
	}

	return nil
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE vote_app SET votes = votes + 1 WHERE category = $1", vote.Category)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getVotesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT category, votes FROM vote_app")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	votes := make(map[string]int)
	for rows.Next() {
		var vote Vote
		if err := rows.Scan(&vote.Category, &vote.Votes); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		votes[vote.Category] = vote.Votes
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}
