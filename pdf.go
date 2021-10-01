package main

import (
	"bytes"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"html/template"
	"io"
	"log"
	"net/http"
)

func GeneratePdf(w http.ResponseWriter, generator ExampleGenerator) error {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("error creating PDF generator: %w", err)
	}
	pdfg.SetOutput(w)
	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Grayscale.Set(true)

	page := wkhtmltopdf.NewPageReader(generator.Generate())
	page.MinimumFontSize.Set(14)
	page.Encoding.Set("utf-8")
	pdfg.AddPage(page)
	err = pdfg.Create()
	if err != nil {
		return fmt.Errorf("error creating PDF: %w", err)
	}
	//err = pdfg.WriteFile(filename)
	//if err != nil {
	//	log.Fatalf("error writing PDF file: %v", err)
	//}
	fmt.Println("PDF has been successfully generated")
	return nil
}

func ExecuteTemplate(text string, data interface{}) io.Reader {
	var buf bytes.Buffer
	t := template.Must(template.New("Template").Parse(text))
	err := t.Execute(&buf, data)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}
	return &buf
}

type ExampleGenerator interface {
	Generate() io.Reader
}
