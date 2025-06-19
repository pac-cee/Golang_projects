package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// createDirectory creates a new directory and returns an error if it fails
func createDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// writeFileBuffered writes data to a file using buffered I/O
func writeFileBuffered(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range data {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to buffer: %w", err)
		}
	}

	return writer.Flush()
}

// readFileBuffered reads a file line by line using buffered I/O
func readFileBuffered(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return lines, nil
}

// getFileInfo retrieves and displays file information
func getFileInfo(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	fmt.Printf("File Information for %s:\n", filename)
	fmt.Printf("Size: %d bytes\n", info.Size())
	fmt.Printf("Permissions: %s\n", info.Mode())
	fmt.Printf("Last Modified: %s\n", info.ModTime())
	fmt.Printf("Is Directory: %t\n", info.IsDir())

	return nil
}

func main() {
	// Create a directory for our files
	dataDir := "data"
	if err := createDirectory(dataDir); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Write data using direct file I/O
	filename := filepath.Join(dataDir, "direct.txt")
	data := []byte("Hello, Direct File I/O!\nThis is written without buffering.")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error writing file directly: %v\n", err)
		return
	}

	// Read the file directly
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file directly: %v\n", err)
		return
	}
	fmt.Println("\n=== Direct File I/O Content ===")
	fmt.Println(string(content))

	// Write data using buffered I/O
	bufferedFile := filepath.Join(dataDir, "buffered.txt")
	lines := []string{
		"Hello, Buffered File I/O!",
		"This is line 2 written with buffering",
		"This is line 3 written with buffering",
		fmt.Sprintf("This line was written at: %s", time.Now().Format(time.RFC3339)),
	}

	if err := writeFileBuffered(bufferedFile, lines); err != nil {
		fmt.Printf("Error writing buffered file: %v\n", err)
		return
	}

	// Read the buffered file
	readLines, err := readFileBuffered(bufferedFile)
	if err != nil {
		fmt.Printf("Error reading buffered file: %v\n", err)
		return
	}

	fmt.Println("\n=== Buffered File I/O Content ===")
	for i, line := range readLines {
		fmt.Printf("Line %d: %s\n", i+1, line)
	}

	// Get and display file information
	fmt.Println("\n=== File Information ===")
	if err := getFileInfo(bufferedFile); err != nil {
		fmt.Printf("Error getting file info: %v\n", err)
		return
	}

	// List files in the data directory
	fmt.Println("\n=== Directory Contents ===")
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("Error getting info for %s: %v\n", entry.Name(), err)
			continue
		}
		fmt.Printf("%s - Size: %d bytes\n", entry.Name(), info.Size())
	}
}
