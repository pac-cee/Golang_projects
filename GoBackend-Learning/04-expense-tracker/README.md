# Project 4: Expense Tracker API (with PostgreSQL)

## Overview
This project demonstrates how to build a Go backend API that connects to a PostgreSQL database. You'll learn how to use Go's `database/sql` package to perform queries and return data as JSON.

## Learning Goals
- Use `database/sql` to interact with relational databases
- Connect to PostgreSQL from Go
- Map SQL query results to Go structs
- Serve data over HTTP as JSON

## Endpoints
- `GET /expenses` — List all expenses from the database

## How to Run
1. Make sure PostgreSQL is running and accessible.
2. Create a database named `expenses` and a table:
   ```sql
   CREATE TABLE expenses (
     id SERIAL PRIMARY KEY,
     title TEXT NOT NULL,
     amount NUMERIC NOT NULL
   );
   INSERT INTO expenses (title, amount) VALUES ('Lunch', 12.5), ('Books', 30.0);
   ```
3. Update the connection string in `main.go` if your DB user/password differ.
4. Run the server:
   ```sh
   go run main.go
   ```

## Why Go?
- Go's standard library provides robust, performant database access
- Easy mapping between SQL rows and Go structs
- Great for building RESTful APIs on top of databases

## Example Usage
- `GET http://localhost:8083/expenses` returns:
  ```json
  [
    {"id":1, "title":"Lunch", "amount":12.5},
    {"id":2, "title":"Books", "amount":30.0}
  ]
  ```

## Next Steps
- Add endpoints for creating, updating, and deleting expenses
- Add user authentication
- Use an ORM like GORM for more advanced mapping

---

This project shows how Go excels at building database-driven backend APIs!
