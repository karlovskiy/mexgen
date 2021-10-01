package main

import (
	"fmt"
	"io"
	"math/rand"
)

const (
	TemplateAddSub = `
<!DOCTYPE html>
<html lang="en">
<body>
<table cellpadding="10" border="1" width="100%">
    {{ range .Rows }}
        <tr>
			{{ range .Columns }}
				<td>
					<table cellpadding="0" cellspacing="0" width="100%">
						<tbody>
							<tr>
								<td rowspan=2><b>{{ .Sign }}</b></td>
								<td style="text-align: right;"><b>{{ .UpperOperand }}</b></td>
							</tr>
							<tr>
								<td style="text-align: right;"><b>{{ .BottomOperand }}</b></td>
							</tr>
							<tr>
								<td>&nbsp;</td>
								<td>&nbsp;</td>
							</tr>
						</tbody>
					</table>
				</td>
			{{ end}}
        </tr>
    {{ end}}
</table>
</body>
</html>
`
)

type AddSub struct {
	Rows int
	Cols int
}

func (e *AddSub) Generate() io.Reader {
	type Column struct {
		Sign          string
		UpperOperand  string
		BottomOperand string
	}
	type Row struct {
		Columns []Column
	}
	type Data struct {
		Rows []Row
	}
	generate := func() (string, string, string) {
		signCode := rand.Intn(2)
		var upper, bottom int
		var sign string
		if signCode == 0 {
			sign = "+"
			upper = rand.Intn(98) + 1        // [1, 98]
			bottom = rand.Intn(99-upper) + 1 // [1, 99-upper]
		} else {
			sign = "-"
			upper = rand.Intn(99) + 1  // [1, 99]
			bottom = rand.Intn(99) + 1 // [1, 99]
		}
		if upper < bottom {
			upper, bottom = bottom, upper
		}
		return sign, fmt.Sprintf("%2d", upper), fmt.Sprintf("%2d", bottom)
	}
	data := Data{}
	for i := 0; i < e.Rows; i++ {
		row := Row{}
		for j := 0; j < e.Cols; j++ {
			sign, upper, bottom := generate()
			col := Column{
				Sign:          sign,
				UpperOperand:  upper,
				BottomOperand: bottom,
			}
			row.Columns = append(row.Columns, col)
		}
		data.Rows = append(data.Rows, row)
	}
	return ExecuteTemplate(TemplateAddSub, data)
}
