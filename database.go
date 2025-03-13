package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
}

func initDatabase() *sql.DB {
	// Use environment variable or default to a path in the /app/data directory
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "/app/data/tasks.db"
	}

	// Ensure the directory exists
	err := os.MkdirAll(filepath.Dir(dbPath), 0755)
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tasks table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'pending'
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createTask(db *sql.DB, title, description string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO tasks (title, description) VALUES (?, ?)",
		title, description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func getTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, title, description, status FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func updateTaskStatus(db *sql.DB, id int, status string) error {
	// VULNERABLE: Direct string concatenation creates SQL injection risk
	query := "UPDATE tasks SET status = '" + status + "' WHERE id = " + fmt.Sprintf("%d", id)
	_, err := db.Exec(query)
	return err
}

func deleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
