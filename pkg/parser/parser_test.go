package parser

import (
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	tests := []struct {
		name     string
		fileType string
		wantType string
		wantErr  bool
	}{
		{"JSON Parser", "json", "*parser.JSONParser", false},
		{"CSV Parser", "csv", "*parser.CSVParser", false},
		{"XML Parser", "xml", "*parser.XMLParser", false},
		{"HTML Parser", "html", "*parser.HTMLParser", false},
		{"Excel Parser", "xlsx", "*parser.ExcelParser", false},
		{"Invalid Parser", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewParser(tt.fileType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got).String() != tt.wantType {
				t.Errorf("NewParser() = %v, want %v", reflect.TypeOf(got), tt.wantType)
			}
		})
	}
}
