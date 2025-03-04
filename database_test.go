package main

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestDatabase() *sql.DB {
	// Use a separate test database
	db, err := sql.Open("sqlite", "./test_tasks.db")
	if err != nil {
		panic(err)
	}

	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'pending'
		)
	`)
	if err != nil {
		panic(err)
	}

	return db
}

func cleanupTestDatabase(db *sql.DB) {
	db.Close()
	os.Remove("./test_tasks.db")
}

func TestCreateTask(t *testing.T) {
	db := setupTestDatabase()
	defer cleanupTestDatabase(db)

	// Test creating a task
	id, err := createTask(db, "Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if id <= 0 {
		t.Errorf("Expected positive task ID, got %d", id)
	}

	// Verify task was created
	rows, err := db.Query("SELECT title, description, status FROM tasks WHERE id = ?", id)
	if err != nil {
		t.Fatalf("Failed to query created task: %v", err)
	}
	defer rows.Close()

	var title, description, status string
	if rows.Next() {
		err = rows.Scan(&title, &description, &status)
		if err != nil {
			t.Fatalf("Failed to scan task: %v", err)
		}

		if title != "Test Task" {
			t.Errorf("Expected title 'Test Task', got %s", title)
		}
		if description != "Test Description" {
			t.Errorf("Expected description 'Test Description', got %s", description)
		}
		if status != "pending" {
			t.Errorf("Expected status 'pending', got %s", status)
		}
	} else {
		t.Error("No task found after creation")
	}
}

func TestGetTasks(t *testing.T) {
	db := setupTestDatabase()
	defer cleanupTestDatabase(db)

	// Create some test tasks
	_, err := createTask(db, "Task 1", "Description 1")
	if err != nil {
		t.Fatalf("Failed to create task 1: %v", err)
	}
	_, err = createTask(db, "Task 2", "Description 2")
	if err != nil {
		t.Fatalf("Failed to create task 2: %v", err)
	}

	// Retrieve tasks
	tasks, err := getTasks(db)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Verify task details
	if tasks[0].Title != "Task 1" || tasks[0].Description != "Description 1" {
		t.Errorf("First task details incorrect: %+v", tasks[0])
	}
	if tasks[1].Title != "Task 2" || tasks[1].Description != "Description 2" {
		t.Errorf("Second task details incorrect: %+v", tasks[1])
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	db := setupTestDatabase()
	defer cleanupTestDatabase(db)

	// Create a task
	id, err := createTask(db, "Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Update task status
	err = updateTaskStatus(db, int(id), "completed")
	if err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	// Verify status update
	rows, err := db.Query("SELECT status FROM tasks WHERE id = ?", id)
	if err != nil {
		t.Fatalf("Failed to query task status: %v", err)
	}
	defer rows.Close()

	var status string
	if rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			t.Fatalf("Failed to scan status: %v", err)
		}

		if status != "completed" {
			t.Errorf("Expected status 'completed', got %s", status)
		}
	} else {
		t.Error("No task found after status update")
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDatabase()
	defer cleanupTestDatabase(db)

	// Create a task
	id, err := createTask(db, "Task to Delete", "Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Delete the task
	err = deleteTask(db, int(id))
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify task was deleted
	rows, err := db.Query("SELECT COUNT(*) FROM tasks WHERE id = ?", id)
	if err != nil {
		t.Fatalf("Failed to query tasks after deletion: %v", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			t.Fatalf("Failed to scan count: %v", err)
		}

		if count != 0 {
			t.Errorf("Expected 0 tasks after deletion, got %d", count)
		}
	}
}
