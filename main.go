package main

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"log"
	"math/rand"
	"time"
)

const (
	ExampleRows    = 32
	ExampleColumns = 8

	PageTemplate = `
<!DOCTYPE html>
<html lang="en">
<body>
<table cellpadding="10" border="1" width="100%">
    {{ range .Rows }}
        <tr>
			{{ range .Columns }}
            	<td><b>{{ .Example }}</b></td>
			{{ end}}
        </tr>
    {{ end}}
</table>
</body>
</html>
`
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalf("error creating PDF generator: %v", err)
	}
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(true)
	data := Data{}
	for i := 0; i < ExampleRows; i++ {
		row := Row{}
		for j := 0; j < ExampleColumns; j++ {
			col := Column{Example: generateRandomExample()}
			row.Columns = append(row.Columns, col)
		}
		data.Rows = append(data.Rows, row)
	}
	t := template.Must(template.New("examplesTemplate").Parse(PageTemplate))
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	page := wkhtmltopdf.NewPageReader(&buf)
	page.MinimumFontSize.Set(14)
	page.Encoding.Set("utf-8")
	pdfg.AddPage(page)
	err = pdfg.Create()
	if err != nil {
		log.Fatalf("error creating PDF: %v", err)
	}
	err = pdfg.WriteFile("./examples.pdf")
	if err != nil {
		log.Fatalf("error writing PDF file: %v", err)
	}
	fmt.Println("Done")
}

func generateRandomExample() string {
	left := rand.Intn(10) + 1
	right := rand.Intn(10) + 1
	operation := rand.Intn(2)
	var operationSign string
	if operation == 0 {
		operationSign = "\u00D7" // multiplication
	} else if operation == 1 {
		operationSign = "\u00F7" // division
		left = left * right
	}
	return fmt.Sprintf("%d%s%d=", left, operationSign, right)
}

type Data struct {
	Rows []Row
}

type Row struct {
	Columns []Column
}

type Column struct {
	Example string
}
