package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
)

type Record struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // Example: Write CSV
    file, _ := os.Create("data.csv")
    writer := csv.NewWriter(file)
    writer.Write([]string{"Name", "Age"})
    writer.Write([]string{"Alice", "30"})
    writer.Flush()
    file.Close()

    // Example: Read CSV and convert to JSON
    file, _ = os.Open("data.csv")
    reader := csv.NewReader(file)
    records, _ := reader.ReadAll()
    file.Close()
    var recs []Record
    for i, row := range records {
        if i == 0 { continue }
        recs = append(recs, Record{Name: row[0]})
    }
    jsonData, _ := json.MarshalIndent(recs, "", "  ")
    fmt.Println(string(jsonData))
}
