package renderer

import (
	"strings"
	"testing"

	"github.com/gowtham2003/gotable/pkg/parser"
)

func TestHTMLRenderer_Render(t *testing.T) {
	testData := &parser.TableData{
		Headers: []string{"Name", "Age"},
		Rows: []map[string]string{
			{"Name": "John", "Age": "30"},
			{"Name": "Alice", "Age": "25"},
		},
	}

	tests := []struct {
		name    string
		data    *parser.TableData
		checks  []string
		wantErr bool
	}{
		{
			name: "Valid HTML Table",
			data: testData,
			checks: []string{
				"<table>",
				"<th>Name</th>",
				"<th>Age</th>",
				"<td>John</td>",
				"<td>30</td>",
				"<td>Alice</td>",
				"<td>25</td>",
				"</table>",
			},
			wantErr: false,
		},
	}

	renderer := NewHTMLRenderer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderer.Render(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("HTMLRenderer.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, check := range tt.checks {
				if !strings.Contains(got, check) {
					t.Errorf("HTMLRenderer.Render() output doesn't contain %q", check)
				}
			}
		})
	}
}
