package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// Print outputs data in the specified format
func Print(data interface{}, format string) {
	switch strings.ToLower(format) {
	case "json":
		printJSON(data)
	case "yaml":
		printYAML(data)
	default:
		printTable(data)
	}
}

func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
	}
}

func printYAML(data interface{}) {
	encoder := yaml.NewEncoder(os.Stdout)
	if err := encoder.Encode(data); err != nil {
		fmt.Printf("Error encoding YAML: %v\n", err)
	}
}

func printTable(data interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Slice {
		fmt.Printf("Error: data must be a slice\n")
		return
	}

	if v.Len() == 0 {
		fmt.Println("No data to display")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Get field names from the first item
	firstItem := v.Index(0)
	if firstItem.Kind() == reflect.Ptr {
		firstItem = firstItem.Elem()
	}
	t := firstItem.Type()

	// Print header
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() {
			headers = append(headers, strings.ToUpper(field.Name))
		}
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Print data rows
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		var row []string
		for j := 0; j < item.NumField(); j++ {
			field := item.Field(j)
			if t.Field(j).IsExported() {
				row = append(row, fmt.Sprintf("%v", field.Interface()))
			}
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
}
