package renderer

import (
	"bytes"
	"testing"

	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/xuri/excelize/v2"
)

func TestExcelRenderer_Render(t *testing.T) {
	testData := &parser.TableData{
		Headers: []string{"Name", "Age"},
		Rows: []map[string]string{
			{"Name": "John", "Age": "30"},
			{"Name": "Alice", "Age": "25"},
		},
	}

	renderer := &ExcelRenderer{}
	output, err := renderer.Render(testData)
	if err != nil {
		t.Fatalf("ExcelRenderer.Render() error = %v", err)
	}

	// Read the generated Excel file
	f, err := excelize.OpenReader(bytes.NewReader([]byte(output)))
	if err != nil {
		t.Fatalf("Failed to read generated Excel: %v", err)
	}
	defer f.Close()

	// Check headers
	for i, header := range testData.Headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		value, _ := f.GetCellValue("Sheet1", cell)
		if value != header {
			t.Errorf("Header mismatch at %s: got %s, want %s", cell, value, header)
		}
	}

	// Check data rows
	for rowIdx, row := range testData.Rows {
		for colIdx, header := range testData.Headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			value, _ := f.GetCellValue("Sheet1", cell)
			if value != row[header] {
				t.Errorf("Data mismatch at %s: got %s, want %s", cell, value, row[header])
			}
		}
	}
}
