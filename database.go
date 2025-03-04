package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID          int
	Title       string
	Description string
	Status      string
}

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "./tasks.db")
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
	_, err := db.Exec(
		"UPDATE tasks SET status = ? WHERE id = ?",
		status, id,
	)
	return err
}

func deleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}
