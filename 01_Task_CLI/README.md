# Task CLI - Go Command Line Task Manager

This is a beginner-friendly Go project that implements a command-line task management system. It demonstrates fundamental Go concepts and file handling.

## Concepts Covered
- Basic Go syntax and data types
- Structs and methods
- File I/O operations
- JSON encoding/decoding
- Error handling
- Command-line interface
- Time handling
- Slices and basic data structures

## Features
- Add new tasks with title and description
- List all tasks
- Update task status
- Delete tasks
- Persistent storage using JSON file
- Simple command-line interface

## Project Structure
```
01_Task_CLI/
├── main.go    # Main application code
├── go.mod     # Go module file
└── tasks.json # Data storage file (created on first run)
```

## How to Run
1. Navigate to the project directory:
   ```bash
   cd 01_Task_CLI
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

## Usage
The application provides a simple menu-driven interface:
1. Add Task - Create a new task
2. List Tasks - View all tasks
3. Update Task Status - Change task status
4. Delete Task - Remove a task
5. Exit - Close the application

## Learning Points
1. **Go Basics**
   - Package declaration
   - Imports
   - Main function
   - Variables and types
   - Control structures

2. **Structs and Methods**
   - Custom type definitions
   - Struct methods
   - Pointer receivers

3. **File Operations**
   - Reading from files
   - Writing to files
   - Error handling

4. **JSON Handling**
   - Marshaling
   - Unmarshaling
   - Struct tags

5. **User Input**
   - Reading from stdin
   - Basic input validation

## Next Steps
- Add task categories
- Implement due dates
- Add search functionality
- Add task priority levels
- Implement data validation
- Add unit tests

This project serves as a foundation for understanding Go basics. The next project will build upon these concepts and introduce web development with Go.
