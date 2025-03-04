package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	db = initDatabase()
	defer db.Close()

	http.HandleFunc("/", handleTaskList)
	http.HandleFunc("/create", handleTaskCreate)
	http.HandleFunc("/toggle", handleTaskToggle)
	http.HandleFunc("/delete", handleTaskDelete)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
