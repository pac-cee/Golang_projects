package main

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "test.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    sqlStmt := `CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, name TEXT);`
    _, err = db.Exec(sqlStmt)
    if err != nil {
        panic(err)
    }
    fmt.Println("Table created!")
}
