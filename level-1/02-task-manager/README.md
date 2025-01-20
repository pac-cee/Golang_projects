# Task Manager

A command-line task management application that demonstrates working with structs, file operations, and JSON in Go.

## Concepts Covered

- Structs and methods
- JSON encoding/decoding
- File I/O operations
- Slices and arrays
- Error handling
- Time handling
- CRUD operations

## Features

- Create new tasks with title, description, and due date
- List all tasks
- Update task status
- Delete tasks
- Persistent storage using JSON file
- Basic error handling

## How to Run

1. Make sure you have Go installed on your system
2. Navigate to the project directory:
   ```bash
   cd level-1/02-task-manager
   ```
3. Run the program:
   ```bash
   go run main.go
   ```

## Usage Example

```
Task Manager
1. Add Task
2. List Tasks
3. Update Task Status
4. Delete Task
5. Exit
Choose option (1-5): 1
Enter task title: Complete Project
Enter task description: Finish the Go learning project
Enter due date (YYYY-MM-DD): 2024-12-31
Task added successfully!
```

## Project Structure

```
02-task-manager/
├── main.go      # Main program file
├── tasks.json   # Data storage file (created when tasks are added)
└── README.md    # Project documentation
```

## Data Structure

Tasks are stored in JSON format with the following structure:
```json
{
  "id": 1,
  "title": "Complete Project",
  "description": "Finish the Go learning project",
  "due_date": "2024-12-31T00:00:00Z",
  "status": "pending"
}
```

## Learning Objectives

- Working with custom types and structs
- Implementing CRUD operations
- File handling in Go
- JSON marshaling and unmarshaling
- Time formatting and parsing
- Error handling patterns
- Using methods with receiver types

## Next Steps

To extend this project, you could:
1. Add task categories or tags
2. Implement task priority levels
3. Add due date notifications
4. Create a basic web interface
5. Add user authentication
6. Implement task filtering and sorting
7. Add unit tests
