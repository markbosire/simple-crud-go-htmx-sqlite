package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func setupTestHandlerDatabase() *sql.DB {
	// Use a separate test database for handlers
	db, err := sql.Open("sqlite", "./test_handlers.db")
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

func cleanupTestHandlerDatabase(db *sql.DB) {
	db.Close()
	os.Remove("./test_handlers.db")
}

func TestHandleTaskCreate(t *testing.T) {
	// Temporarily replace global db with test db
	originalDB := db
	testDB := setupTestHandlerDatabase()
	db = testDB
	defer func() {
		db = originalDB
		cleanupTestHandlerDatabase(testDB)
	}()

	// Test POST request to create task
	form := url.Values{}
	form.Add("title", "Test Task")
	form.Add("description", "Test Description")

	req, err := http.NewRequest("POST", "/create", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handleTaskCreate(w, req)

	// Check response status code
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, w.Code)
	}

	// Verify task was created in database
	tasks, err := getTasks(testDB)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	if tasks[0].Title != "Test Task" || tasks[0].Description != "Test Description" {
		t.Errorf("Task details incorrect: %+v", tasks[0])
	}

	// Success message
	println("✅ TestHandleTaskCreate: Task creation handler test passed successfully!")
}

func TestHandleTaskToggle(t *testing.T) {
	// Temporarily replace global db with test db
	originalDB := db
	testDB := setupTestHandlerDatabase()
	db = testDB
	defer func() {
		db = originalDB
		cleanupTestHandlerDatabase(testDB)
	}()

	// Create a task first
	taskID, err := createTask(testDB, "Toggle Task", "Toggle Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test toggle to completed
	req, err := http.NewRequest("GET", "/toggle?id="+strconv.Itoa(int(taskID)), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handleTaskToggle(w, req)

	// Check response status code
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, w.Code)
	}

	// Verify task status was updated
	tasks, err := getTasks(testDB)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if tasks[0].Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", tasks[0].Status)
	}

	// Test toggle back to pending
	req, err = http.NewRequest("GET", "/toggle?id="+strconv.Itoa(int(taskID)), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w = httptest.NewRecorder()
	handleTaskToggle(w, req)

	// Verify task status was updated again
	tasks, err = getTasks(testDB)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if tasks[0].Status != "pending" {
		t.Errorf("Expected status 'pending', got %s", tasks[0].Status)
	}

	// Success message
	println("✅ TestHandleTaskToggle: Task toggle handler test passed successfully!")
}

func TestHandleTaskDelete(t *testing.T) {
	// Temporarily replace global db with test db
	originalDB := db
	testDB := setupTestHandlerDatabase()
	db = testDB
	defer func() {
		db = originalDB
		cleanupTestHandlerDatabase(testDB)
	}()

	// Create a task first
	taskID, err := createTask(testDB, "Delete Task", "Delete Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test delete task
	req, err := http.NewRequest("GET", "/delete?id="+strconv.Itoa(int(taskID)), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	handleTaskDelete(w, req)

	// Check response status code
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, w.Code)
	}

	// Verify task was deleted
	tasks, err := getTasks(testDB)
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	// Success message
	println("✅ TestHandleTaskDelete: Task delete handler test passed successfully!")
}
