package renderer

import (
	"reflect"
	"testing"
)

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		wantType string
		wantErr  bool
	}{
		{"ASCII Renderer", "ascii", "*renderer.ASCIIRenderer", false},
		{"CSV Renderer", "csv", "*renderer.CSVRenderer", false},
		{"JSON Renderer", "json", "*renderer.JSONRenderer", false},
		{"HTML Renderer", "html", "*renderer.HTMLRenderer", false},
		{"Excel Renderer", "xlsx", "*renderer.ExcelRenderer", false},
		{"PNG Renderer", "png", "*renderer.ImageRenderer", false},
		{"Invalid Renderer", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRenderer(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRenderer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && reflect.TypeOf(got).String() != tt.wantType {
				t.Errorf("NewRenderer() = %v, want %v", reflect.TypeOf(got), tt.wantType)
			}
		})
	}
}
