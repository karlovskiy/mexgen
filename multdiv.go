package main

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	TemplateMultDiv = `
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

type MultDiv struct {
	Rows int
	Cols int
}

func (e *MultDiv) Generate() io.Reader {
	type Column struct {
		Example string
	}
	type Row struct {
		Columns []Column
	}
	type Data struct {
		Rows []Row
	}
	generate := func() string {
		left := rand.Intn(10) + 1  // [1, 10]
		right := rand.Intn(10) + 1 // [1, 10]
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
	data := Data{}
	for i := 0; i < e.Rows; i++ {
		row := Row{}
		for j := 0; j < e.Cols; j++ {
			col := Column{Example: generate()}
			row.Columns = append(row.Columns, col)
		}
		data.Rows = append(data.Rows, row)
	}
	return ExecuteTemplate(TemplateMultDiv, data)
}
