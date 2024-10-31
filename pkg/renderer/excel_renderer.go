package renderer

import (
	"bytes"

	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/xuri/excelize/v2"
)

type ExcelRenderer struct{}

func (r *ExcelRenderer) Render(data *parser.TableData) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Get the default sheet name
	sheetName := "Sheet1"

	// Write headers
	for col, header := range data.Headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data rows
	for rowIdx, row := range data.Rows {
		for colIdx, header := range data.Headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, row[header])
		}
	}

	// Apply styling
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#f2f2f2"},
			Pattern: 1,
		},
	})
	if err == nil {
		// Apply style to header row
		f.SetRowStyle(sheetName, 1, 1, style)
	}

	// Auto-fit columns
	for col := 1; col <= len(data.Headers); col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		f.SetColWidth(sheetName, colName, colName, 15)
	}

	// Save to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
