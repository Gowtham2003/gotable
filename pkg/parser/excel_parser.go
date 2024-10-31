package parser

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelParser struct{}

func (p *ExcelParser) Parse(input []byte) (*TableData, error) {
	// Create a temporary file from input bytes
	f, err := excelize.OpenReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get the first sheet name
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Get all rows from the first sheet
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must contain at least headers and one data row")
	}

	// First row as headers
	headers := rows[0]

	// Process data rows
	tableRows := make([]map[string]string, 0, len(rows)-1)
	for _, row := range rows[1:] {
		rowData := make(map[string]string)
		for i, cell := range row {
			if i < len(headers) {
				rowData[headers[i]] = cell
			}
		}
		tableRows = append(tableRows, rowData)
	}

	return &TableData{
		Headers: headers,
		Rows:    tableRows,
	}, nil
}
