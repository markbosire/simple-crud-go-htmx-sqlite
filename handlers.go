package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Create a struct for snackbar messages
type SnackbarMessage struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func setSnackbarHeader(w http.ResponseWriter, message string, messageType string) {
	snackbar := SnackbarMessage{
		Message: message,
		Type:    messageType,
	}

	// Convert snackbar to JSON
	snackbarJSON, err := json.Marshal(map[string]SnackbarMessage{
		"showSnackbar": snackbar,
	})
	if err != nil {
		// Fallback if JSON marshaling fails
		w.Header().Add("HX-Trigger", `{"showSnackbar": {"message": "Operation completed", "type": "success"}}`)
		return
	}

	// Set the header with the JSON-encoded snackbar message
	w.Header().Add("HX-Trigger", string(snackbarJSON))
}
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	// Parse both base and specific template
	t, err := template.ParseFiles(
		"templates/base.html",
		fmt.Sprintf("templates/%s", tmpl),
	)
	if err != nil {
		http.Error(w, "Error parsing templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the base template with the content template
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

func handleTaskList(w http.ResponseWriter, r *http.Request) {
	tasks, err := getTasks(db)
	if err != nil {
		renderTemplate(w, "task-list.html", map[string]interface{}{
			"Tasks": tasks,
			"Error": "Failed to load tasks",
		})
		return
	}
	renderTemplate(w, "task-list.html", tasks)
}
func handleTaskCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		description := r.FormValue("description")

		_, err := createTask(db, title, description)
		if err != nil {
			renderTemplate(w, "task-form.html", map[string]interface{}{
				"Error": "Failed to create task",
			})
			setSnackbarHeader(w, "Failed to create task", "error")
			return
		}

		// Use the new helper function to set snackbar header
		setSnackbarHeader(w, "Task created successfully", "success")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		renderTemplate(w, "task-form.html", nil)
	}
}

// Update other handlers similarly with setSnackbarHeader
func handleTaskToggle(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		setSnackbarHeader(w, "Invalid task ID", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Directly use the first method to toggle status
	currentStatus, err := func() (string, error) {
		var status string
		err := db.QueryRow("SELECT status FROM tasks WHERE id = ?", id).Scan(&status)
		return status, err
	}()

	if err != nil {
		log.Printf("Error retrieving task status: %v", err)
		setSnackbarHeader(w, "Failed to retrieve task status", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Determine new status
	newStatus := "completed"
	if currentStatus == "completed" {
		newStatus = "pending"
	}

	// Update status
	_, err = db.Exec("UPDATE tasks SET status = ? WHERE id = ?", newStatus, id)
	if err != nil {
		log.Printf("Error updating task status: %v", err)
		setSnackbarHeader(w, "Failed to update task status", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setSnackbarHeader(w, "Task status updated", "success")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleTaskDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	err := deleteTask(db, id)
	if err != nil {
		setSnackbarHeader(w, "Failed to delete task", "error")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	setSnackbarHeader(w, "Task deleted successfully", "success")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ... rest of the code remains the same
