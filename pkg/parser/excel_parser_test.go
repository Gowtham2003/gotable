package parser

import (
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"
)

func createTestExcelFile() []byte {
	f := excelize.NewFile()
	defer f.Close()

	// Add headers
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")

	// Add data
	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", "30")
	f.SetCellValue("Sheet1", "A3", "Alice")
	f.SetCellValue("Sheet1", "B3", "25")

	// Save to buffer
	buf, _ := f.WriteToBuffer()
	return buf.Bytes()
}

func TestExcelParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *TableData
		wantErr bool
	}{
		{
			name:  "Valid Excel",
			input: createTestExcelFile(),
			want: &TableData{
				Headers: []string{"Name", "Age"},
				Rows: []map[string]string{
					{"Name": "John", "Age": "30"},
					{"Name": "Alice", "Age": "25"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid Excel",
			input:   []byte("invalid excel"),
			want:    nil,
			wantErr: true,
		},
	}

	parser := &ExcelParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExcelParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExcelParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
