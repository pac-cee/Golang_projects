package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Task represents a single task with its properties
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "pending", "in_progress", "completed"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskManager handles all task-related operations
type TaskManager struct {
	tasks    []Task
	filename string
}

// NewTaskManager creates a new instance of TaskManager
func NewTaskManager(filename string) *TaskManager {
	return &TaskManager{
		tasks:    make([]Task, 0),
		filename: filename,
	}
}

// LoadTasks loads tasks from the JSON file
func (tm *TaskManager) LoadTasks() error {
	file, err := os.ReadFile(tm.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, &tm.tasks)
}

// SaveTasks saves tasks to the JSON file
func (tm *TaskManager) SaveTasks() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling tasks: %v", err)
	}

	return os.WriteFile(tm.filename, data, 0644)
}

// AddTask adds a new task to the list
func (tm *TaskManager) AddTask(title, description string) error {
	task := Task{
		ID:          len(tm.tasks) + 1,
		Title:       title,
		Description: description,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tm.tasks = append(tm.tasks, task)
	return tm.SaveTasks()
}

// ListTasks displays all tasks
func (tm *TaskManager) ListTasks() {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	for _, task := range tm.tasks {
		fmt.Printf("\nID: %d\nTitle: %s\nDescription: %s\nStatus: %s\nCreated: %v\nUpdated: %v\n",
			task.ID, task.Title, task.Description, task.Status, task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("----------------------------------------")
	}
}

// UpdateTaskStatus updates the status of a task
func (tm *TaskManager) UpdateTaskStatus(id int, status string) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Status = status
			tm.tasks[i].UpdatedAt = time.Now()
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// DeleteTask removes a task by ID
func (tm *TaskManager) DeleteTask(id int) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func main() {
	tm := NewTaskManager("tasks.json")
	err := tm.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nTask Manager CLI")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Update Task Status")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")
		fmt.Print("Choose an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter task title: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Enter task description: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			err := tm.AddTask(title, desc)
			if err != nil {
				fmt.Printf("Error adding task: %v\n", err)
			} else {
				fmt.Println("Task added successfully!")
			}

		case "2":
			tm.ListTasks()

		case "3":
			fmt.Print("Enter task ID: ")
			var id int
			fmt.Scanf("%d\n", &id)

			fmt.Print("Enter new status (pending/in_progress/completed): ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)

			err := tm.UpdateTaskStatus(id, status)
			if err != nil {
				fmt.Printf("Error updating task: %v\n", err)
			} else {
				fmt.Println("Task updated successfully!")
			}

		case "4":
			fmt.Print("Enter task ID to delete: ")
			var id int
			fmt.Scanf("%d\n", &id)

			err := tm.DeleteTask(id)
			if err != nil {
				fmt.Printf("Error deleting task: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully!")
			}

		case "5":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}
