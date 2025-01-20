package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
}

type TaskManager struct {
	tasks    []Task
	filename string
}

func NewTaskManager(filename string) *TaskManager {
	return &TaskManager{
		tasks:    make([]Task, 0),
		filename: filename,
	}
}

func (tm *TaskManager) LoadTasks() error {
	data, err := ioutil.ReadFile(tm.filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &tm.tasks)
}

func (tm *TaskManager) SaveTasks() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(tm.filename, data, 0644)
}

func (tm *TaskManager) AddTask(title, description string, dueDate time.Time) {
	task := Task{
		ID:          len(tm.tasks) + 1,
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      "pending",
	}
	tm.tasks = append(tm.tasks, task)
	tm.SaveTasks()
}

func (tm *TaskManager) ListTasks() {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("\nTasks:")
	fmt.Println("----------------------------------------")
	for _, task := range tm.tasks {
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Due Date: %s\n", task.DueDate.Format("2006-01-02"))
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Println("----------------------------------------")
	}
}

func (tm *TaskManager) UpdateTaskStatus(id int, status string) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Status = status
			return tm.SaveTasks()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

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

	for {
		fmt.Println("\nTask Manager")
		fmt.Println("1. Add Task")
		fmt.Println("2. List Tasks")
		fmt.Println("3. Update Task Status")
		fmt.Println("4. Delete Task")
		fmt.Println("5. Exit")
		fmt.Print("Choose option (1-5): ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			var title, description, dateStr string
			fmt.Print("Enter task title: ")
			fmt.Scan(&title)
			fmt.Print("Enter task description: ")
			fmt.Scan(&description)
			fmt.Print("Enter due date (YYYY-MM-DD): ")
			fmt.Scan(&dateStr)

			dueDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				fmt.Println("Invalid date format!")
				continue
			}

			tm.AddTask(title, description, dueDate)
			fmt.Println("Task added successfully!")

		case 2:
			tm.ListTasks()

		case 3:
			var id int
			var status string
			fmt.Print("Enter task ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter new status (pending/completed): ")
			fmt.Scan(&status)

			err := tm.UpdateTaskStatus(id, status)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Task status updated successfully!")
			}

		case 4:
			var id int
			fmt.Print("Enter task ID: ")
			fmt.Scan(&id)

			err := tm.DeleteTask(id)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully!")
			}

		case 5:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid choice!")
		}
	}
}
