package renderer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gowtham2003/gotable/pkg/parser"
)

// Renderer interface for different output formats
type Renderer interface {
	Render(data *parser.TableData) (string, error)
}

// Styleable interface for renderers that support styling
type Styleable interface {
	SetStyle(StyleOptions)
}

// StyleOptions contains configuration for rendering styles
type StyleOptions struct {
	BorderStyle  string
	ColorEnabled bool
	FontFamily   string
	TableWidth   int
	// Add other style options as needed
}

// ASCIIRenderer implements Renderer for ASCII table output
type ASCIIRenderer struct{}

// CSVRenderer implements Renderer for CSV output
type CSVRenderer struct{}

// JSONRenderer implements Renderer for JSON output
type JSONRenderer struct{}

// MarkdownRenderer implements Renderer for Markdown table output
type MarkdownRenderer struct{}

func NewRenderer(format string) (Renderer, error) {
	switch strings.ToLower(format) {
	case "ascii":
		return &ASCIIRenderer{}, nil
	case "csv":
		return &CSVRenderer{}, nil
	case "json":
		return &JSONRenderer{}, nil
	case "markdown":
		return &MarkdownRenderer{}, nil
	case "png":
		return NewImageRenderer(), nil
	case "html":
		return NewHTMLRenderer(), nil
	case "xlsx":
		return &ExcelRenderer{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func (r *ASCIIRenderer) Render(data *parser.TableData) (string, error) {
	var result strings.Builder
	widths := getColumnWidths(data)

	// Create separator line
	separator := createSeparator(data.Headers, widths)
	result.WriteString(separator)

	// Write headers
	result.WriteString("|")
	for _, h := range data.Headers {
		format := fmt.Sprintf(" %%-%ds |", widths[h])
		result.WriteString(fmt.Sprintf(format, h))
	}
	result.WriteString("\n")
	result.WriteString(separator)

	// Write data rows
	for _, row := range data.Rows {
		result.WriteString("|")
		for _, h := range data.Headers {
			format := fmt.Sprintf(" %%-%ds |", widths[h])
			result.WriteString(fmt.Sprintf(format, row[h]))
		}
		result.WriteString("\n")
	}
	result.WriteString(separator)

	return result.String(), nil
}

func (r *CSVRenderer) Render(data *parser.TableData) (string, error) {
	var result strings.Builder
	writer := csv.NewWriter(&result)

	// Write headers
	if err := writer.Write(data.Headers); err != nil {
		return "", err
	}

	// Write rows
	for _, row := range data.Rows {
		record := make([]string, len(data.Headers))
		for i, h := range data.Headers {
			record[i] = row[h]
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	return result.String(), writer.Error()
}

func (r *JSONRenderer) Render(data *parser.TableData) (string, error) {
	output := make([]map[string]string, len(data.Rows)+1)

	// Add headers as first row
	headerRow := make(map[string]string)
	for _, h := range data.Headers {
		headerRow[h] = h
	}
	output[0] = headerRow

	// Add data rows
	copy(output[1:], data.Rows)

	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (r *MarkdownRenderer) Render(data *parser.TableData) (string, error) {
	var result strings.Builder

	// Write headers
	result.WriteString("| ")
	result.WriteString(strings.Join(data.Headers, " | "))
	result.WriteString(" |\n")

	// Write separator
	result.WriteString("| ")
	for range data.Headers {
		result.WriteString("--- | ")
	}
	result.WriteString("\n")

	// Write rows
	for _, row := range data.Rows {
		result.WriteString("| ")
		for _, h := range data.Headers {
			result.WriteString(row[h] + " | ")
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

// Helper functions
func getColumnWidths(data *parser.TableData) map[string]int {
	widths := make(map[string]int)

	// Initialize with header lengths
	for _, h := range data.Headers {
		widths[h] = len(h)
	}

	// Check all rows for maximum width
	for _, row := range data.Rows {
		for _, h := range data.Headers {
			if width := len(row[h]); width > widths[h] {
				widths[h] = width
			}
		}
	}

	return widths
}

func createSeparator(headers []string, widths map[string]int) string {
	var sep strings.Builder
	sep.WriteString("+")
	for _, h := range headers {
		sep.WriteString(strings.Repeat("-", widths[h]+2))
		sep.WriteString("+")
	}
	sep.WriteString("\n")
	return sep.String()
}
