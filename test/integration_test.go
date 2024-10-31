package test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/gowtham2003/gotable/pkg/renderer"
	"github.com/xuri/excelize/v2"
	"golang.org/x/net/html"
)

// Test data structures
var (
	simpleData = `[
		{"name": "Name", "age": "Age", "city": "City"},
		{"name": "John", "age": "30", "city": "New York"},
		{"name": "Alice", "age": "25", "city": "London"}
	]`

	complexData = `[
		{"id": "ID", "name": "Name", "age": "Age", "city": "City", "salary": "Salary", "department": "Department"},
		{"id": "1", "name": "John Doe", "age": "30", "city": "New York", "salary": "75000", "department": "Engineering"},
		{"id": "2", "name": "Alice Smith", "age": "25", "city": "London", "salary": "65000", "department": "Marketing"},
		{"id": "3", "name": "Bob Johnson", "age": "35", "city": "Paris", "salary": "80000", "department": "Engineering"},
		{"id": "4", "name": "Carol White", "age": "28", "city": "Tokyo", "salary": "70000", "department": "Sales"}
	]`

	emptyData = `[
		{"name": "Name", "age": "Age", "city": "City"}
	]`
)

// Helper functions for validation
func validateJSONOutput(t *testing.T, data []byte) bool {
	var result []map[string]string
	return json.Unmarshal(data, &result) == nil
}

func validateHTMLOutput(t *testing.T, data []byte) bool {
	_, err := html.Parse(strings.NewReader(string(data)))
	return err == nil
}

func validateExcelOutput(t *testing.T, data []byte) bool {
	_, err := excelize.OpenReader(strings.NewReader(string(data)))
	return err == nil
}

func validateCSVOutput(t *testing.T, data []byte) bool {
	return strings.Count(string(data), "\n") > 0 && strings.Contains(string(data), ",")
}

func validateMarkdownOutput(t *testing.T, data []byte) bool {
	content := string(data)
	return strings.Contains(content, "|") && strings.Contains(content, "---")
}

func validateASCIIOutput(t *testing.T, data []byte) bool {
	content := string(data)
	return strings.Contains(content, "+") && strings.Contains(content, "|")
}

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name         string
		inputFormat  string
		outputFormat string
		inputData    string
		validator    func(*testing.T, []byte) bool
		description  string
	}{
		// JSON input to various outputs
		{
			name:         "JSON to JSON",
			inputFormat:  "json",
			outputFormat: "json",
			inputData:    simpleData,
			validator:    validateJSONOutput,
			description:  "Convert JSON to JSON (roundtrip)",
		},
		{
			name:         "JSON to HTML",
			inputFormat:  "json",
			outputFormat: "html",
			inputData:    simpleData,
			validator:    validateHTMLOutput,
			description:  "Convert JSON to HTML table",
		},
		{
			name:         "JSON to Excel",
			inputFormat:  "json",
			outputFormat: "xlsx",
			inputData:    simpleData,
			validator:    validateExcelOutput,
			description:  "Convert JSON to Excel spreadsheet",
		},
		{
			name:         "JSON to CSV",
			inputFormat:  "json",
			outputFormat: "csv",
			inputData:    simpleData,
			validator:    validateCSVOutput,
			description:  "Convert JSON to CSV format",
		},
		{
			name:         "JSON to Markdown",
			inputFormat:  "json",
			outputFormat: "markdown",
			inputData:    simpleData,
			validator:    validateMarkdownOutput,
			description:  "Convert JSON to Markdown table",
		},
		{
			name:         "JSON to ASCII",
			inputFormat:  "json",
			outputFormat: "ascii",
			inputData:    simpleData,
			validator:    validateASCIIOutput,
			description:  "Convert JSON to ASCII table",
		},

		// Complex data tests
		{
			name:         "Complex JSON to HTML",
			inputFormat:  "json",
			outputFormat: "html",
			inputData:    complexData,
			validator:    validateHTMLOutput,
			description:  "Convert complex JSON data to HTML table",
		},
		{
			name:         "Complex JSON to Excel",
			inputFormat:  "json",
			outputFormat: "xlsx",
			inputData:    complexData,
			validator:    validateExcelOutput,
			description:  "Convert complex JSON data to Excel spreadsheet",
		},

		// Empty data tests
		{
			name:         "Empty JSON to HTML",
			inputFormat:  "json",
			outputFormat: "html",
			inputData:    emptyData,
			validator:    validateHTMLOutput,
			description:  "Convert empty JSON data to HTML table",
		},

		// CSV input to various outputs
		{
			name:         "CSV to JSON",
			inputFormat:  "csv",
			outputFormat: "json",
			inputData:    "name,age,city\nJohn,30,New York\nAlice,25,London",
			validator:    validateJSONOutput,
			description:  "Convert CSV to JSON format",
		},
		{
			name:         "CSV to HTML",
			inputFormat:  "csv",
			outputFormat: "html",
			inputData:    "name,age,city\nJohn,30,New York\nAlice,25,London",
			validator:    validateHTMLOutput,
			description:  "Convert CSV to HTML table",
		},

		// Excel input to various outputs
		{
			name:         "Excel to JSON",
			inputFormat:  "xlsx",
			outputFormat: "json",
			inputData:    createTestExcelData(),
			validator:    validateJSONOutput,
			description:  "Convert Excel to JSON format",
		},
		{
			name:         "Excel to HTML",
			inputFormat:  "xlsx",
			outputFormat: "html",
			inputData:    createTestExcelData(),
			validator:    validateHTMLOutput,
			description:  "Convert Excel to HTML table",
		},

		// HTML input to various outputs
		{
			name:         "HTML to JSON",
			inputFormat:  "html",
			outputFormat: "json",
			inputData:    createTestHTMLTable(),
			validator:    validateJSONOutput,
			description:  "Convert HTML table to JSON format",
		},
		{
			name:         "HTML to Excel",
			inputFormat:  "html",
			outputFormat: "xlsx",
			inputData:    createTestHTMLTable(),
			validator:    validateExcelOutput,
			description:  "Convert HTML table to Excel spreadsheet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test files
			tempDir := t.TempDir()

			// Create input file
			inputFile := filepath.Join(tempDir, "input."+tt.inputFormat)
			if err := os.WriteFile(inputFile, []byte(tt.inputData), 0644); err != nil {
				t.Fatalf("Failed to create input file: %v", err)
			}

			// Create output file path
			outputFile := filepath.Join(tempDir, "output."+tt.outputFormat)

			// Create parser
			p, err := parser.NewParser(tt.inputFormat)
			if err != nil {
				t.Fatalf("Failed to create parser: %v", err)
			}

			// Parse input
			input, err := os.ReadFile(inputFile)
			if err != nil {
				t.Fatalf("Failed to read input file: %v", err)
			}

			data, err := p.Parse(input)
			if err != nil {
				t.Fatalf("Failed to parse input: %v", err)
			}

			// Create renderer
			r, err := renderer.NewRenderer(tt.outputFormat)
			if err != nil {
				t.Fatalf("Failed to create renderer: %v", err)
			}

			// Render output
			output, err := r.Render(data)
			if err != nil {
				t.Fatalf("Failed to render output: %v", err)
			}

			// Write output
			if err := os.WriteFile(outputFile, []byte(output), 0644); err != nil {
				t.Fatalf("Failed to write output file: %v", err)
			}

			// Read and validate output
			outputData, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			if !tt.validator(t, outputData) {
				t.Errorf("Output validation failed for %s", tt.name)
			}

			// Additional checks
			info, err := os.Stat(outputFile)
			if err != nil {
				t.Fatalf("Failed to stat output file: %v", err)
			}
			if info.Size() == 0 {
				t.Error("Output file is empty")
			}
		})
	}
}

// Helper function to create test Excel data
func createTestExcelData() string {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")
	f.SetCellValue("Sheet1", "C1", "City")
	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", "30")
	f.SetCellValue("Sheet1", "C2", "New York")

	var buf strings.Builder
	if err := f.Write(&buf); err != nil {
		return ""
	}
	return buf.String()
}

// Helper function to create test HTML table
func createTestHTMLTable() string {
	return `
		<table>
			<tr>
				<th>Name</th>
				<th>Age</th>
				<th>City</th>
			</tr>
			<tr>
				<td>John</td>
				<td>30</td>
				<td>New York</td>
			</tr>
			<tr>
				<td>Alice</td>
				<td>25</td>
				<td>London</td>
			</tr>
		</table>
	`
}

// Error test cases
func TestEndToEndErrors(t *testing.T) {
	errorTests := []struct {
		name         string
		inputFormat  string
		outputFormat string
		inputData    string
		expectError  bool
	}{
		{
			name:         "Invalid JSON Input",
			inputFormat:  "json",
			outputFormat: "html",
			inputData:    "invalid json",
			expectError:  true,
		},
		{
			name:         "Invalid CSV Input",
			inputFormat:  "csv",
			outputFormat: "json",
			inputData:    "invalid,csv,format\ninvalid",
			expectError:  true,
		},
		{
			name:         "Invalid Format Combination",
			inputFormat:  "invalid",
			outputFormat: "invalid",
			inputData:    "{}",
			expectError:  true,
		},
		{
			name:         "Empty Input",
			inputFormat:  "json",
			outputFormat: "html",
			inputData:    "",
			expectError:  true,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			inputFile := filepath.Join(tempDir, "input."+tt.inputFormat)
			filepath.Join(tempDir, "output."+tt.outputFormat)

			if err := os.WriteFile(inputFile, []byte(tt.inputData), 0644); err != nil {
				t.Fatalf("Failed to create input file: %v", err)
			}

			p, err := parser.NewParser(tt.inputFormat)
			if err != nil && !tt.expectError {
				t.Fatalf("Unexpected error creating parser: %v", err)
			}

			if err == nil {
				input, _ := os.ReadFile(inputFile)
				_, err = p.Parse(input)
				if err == nil && tt.expectError {
					t.Error("Expected error but got none")
				}
			}
		})
	}
}
