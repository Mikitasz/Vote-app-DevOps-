package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Vote struct {
	Category string `json:"category"`
	Votes    int    `json:"votes"`
}

func main() {
	var err error
	connStr := "user=postgres dbname=vote_app sslmode=disable password=mysecret host=localhost"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/votes", getVotesHandler)
	http.Handle("/", http.FileServer(http.Dir("template")))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
