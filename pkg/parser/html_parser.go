package parser

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type HTMLParser struct{}

func (p *HTMLParser) Parse(input []byte) (*TableData, error) {
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	// Find the first table element
	table := findFirstTable(doc)
	if table == nil {
		return nil, fmt.Errorf("no table found in HTML")
	}

	var headers []string
	var rows []map[string]string

	// Process table rows
	var processingHeader bool = true
	for row := findFirstTag(table, "tr"); row != nil; row = findNextSibling(row, "tr") {
		if processingHeader {
			// Process header row
			for cell := findFirstTag(row, "th"); cell != nil; cell = findNextSibling(cell, "th") {
				headers = append(headers, getText(cell))
			}

			// If no th tags found, try td tags for header
			if len(headers) == 0 {
				for cell := findFirstTag(row, "td"); cell != nil; cell = findNextSibling(cell, "td") {
					headers = append(headers, getText(cell))
				}
			}
			processingHeader = false
			continue
		}

		// Process data rows
		rowData := make(map[string]string)
		i := 0
		for cell := findFirstTag(row, "td"); cell != nil; cell = findNextSibling(cell, "td") {
			if i < len(headers) {
				rowData[headers[i]] = getText(cell)
				i++
			}
		}
		if len(rowData) > 0 {
			rows = append(rows, rowData)
		}
	}

	return &TableData{
		Headers: headers,
		Rows:    rows,
	}, nil
}

// Helper functions for HTML parsing
func findFirstTable(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "table" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if result := findFirstTable(c); result != nil {
			return result
		}
	}
	return nil
}

func findFirstTag(n *html.Node, tag string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func findNextSibling(n *html.Node, tag string) *html.Node {
	for c := n.NextSibling; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == tag {
			return c
		}
	}
	return nil
}

func getText(n *html.Node) string {
	var text strings.Builder
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}
	extract(n)
	return strings.TrimSpace(text.String())
}
