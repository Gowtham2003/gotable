package parser

import (
	"reflect"
	"testing"
)

func TestJSONParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *TableData
		wantErr bool
	}{
		{
			name: "Valid JSON",
			input: `[
				{"name": "Name", "age": "Age"},
				{"name": "John", "age": "30"},
				{"name": "Alice", "age": "25"}
			]`,
			want: &TableData{
				Headers: []string{"name", "age"},
				Rows: []map[string]string{
					{"name": "John", "age": "30"},
					{"name": "Alice", "age": "25"},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   `invalid json`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty JSON Array",
			input:   `[]`,
			want:    nil,
			wantErr: true,
		},
	}

	parser := &JSONParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.Parse([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONParser.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
