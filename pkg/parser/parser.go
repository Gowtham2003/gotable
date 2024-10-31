package parser

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// TableData represents the parsed data structure
type TableData struct {
	Headers []string
	Rows    []map[string]string
}

// Parser interface for different input formats
type Parser interface {
	Parse(input []byte) (*TableData, error)
}

// JSONParser implements Parser for JSON input
type JSONParser struct{}

// CSVParser implements Parser for CSV input
type CSVParser struct{}

// XMLParser implements Parser for XML input
type XMLParser struct{}

func NewParser(fileType string) (Parser, error) {
	switch strings.ToLower(fileType) {
	case "json":
		return &JSONParser{}, nil
	case "csv":
		return &CSVParser{}, nil
	case "xml":
		return &XMLParser{}, nil
	case "html":
		return &HTMLParser{}, nil
	case "xlsx":
		return &ExcelParser{}, nil
	default:
		return nil, fmt.Errorf("unsupported input format: %s", fileType)
	}
}

func (p *JSONParser) Parse(input []byte) (*TableData, error) {
	var rawData []map[string]interface{}
	if err := json.Unmarshal(input, &rawData); err != nil {
		return nil, err
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("empty JSON array")
	}

	// Extract headers from first object
	headers := make([]string, 0)
	for key := range rawData[0] {
		headers = append(headers, key)
	}

	// Convert data rows
	rows := make([]map[string]string, len(rawData)-1)
	for i := 1; i < len(rawData); i++ {
		row := make(map[string]string)
		for _, h := range headers {
			row[h] = fmt.Sprintf("%v", rawData[i][h])
		}
		rows[i-1] = row
	}

	return &TableData{
		Headers: headers,
		Rows:    rows,
	}, nil
}

func (p *CSVParser) Parse(input []byte) (*TableData, error) {
	reader := csv.NewReader(strings.NewReader(string(input)))

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Read rows
	var rows []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for i, value := range record {
			row[headers[i]] = value
		}
		rows = append(rows, row)
	}

	return &TableData{
		Headers: headers,
		Rows:    rows,
	}, nil
}

func (p *XMLParser) Parse(input []byte) (*TableData, error) {
	// Implementation depends on your XML structure
	// This is a basic example
	var data struct {
		Rows []map[string]string `xml:"row"`
	}

	if err := xml.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	if len(data.Rows) == 0 {
		return nil, fmt.Errorf("empty XML data")
	}

	// Extract headers from first row
	headers := make([]string, 0)
	for key := range data.Rows[0] {
		headers = append(headers, key)
	}

	return &TableData{
		Headers: headers,
		Rows:    data.Rows[1:],
	}, nil
}
