package main

import (
	"fmt"
	"os"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input.html> [output.pdf]")
		return
	}
	input := os.Args[1]
	output := "output.pdf"
	if len(os.Args) > 2 {
		output = os.Args[2]
	}
	f, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening HTML file:", err)
		return
	}
	defer f.Close()

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println("Error creating PDF generator:", err)
		return
	}
	page := wkhtmltopdf.NewPageReader(f)
	pdfg.AddPage(page)
	if err := pdfg.Create(); err != nil {
		fmt.Println("Error creating PDF:", err)
		return
	}
	if err := pdfg.WriteFile(output); err != nil {
		fmt.Println("Error writing PDF:", err)
		return
	}
	fmt.Println("PDF generated:", output)
}
