package renderer

import (
	"fmt"
	"strings"

	"github.com/gowtham2003/gotable/pkg/parser"
)

type HTMLRenderer struct {
	Style string
}

func NewHTMLRenderer() *HTMLRenderer {
	return &HTMLRenderer{
		Style: `
			<style>
				table {
					border-collapse: collapse;
					width: 100%;
					margin: 20px 0;
					font-family: Arial, sans-serif;
				}
				th, td {
					border: 1px solid #ddd;
					padding: 8px;
					text-align: left;
				}
				th {
					background-color: #f2f2f2;
					font-weight: bold;
				}
				tr:nth-child(even) {
					background-color: #f9f9f9;
				}
				tr:hover {
					background-color: #f5f5f5;
				}
			</style>
		`,
	}
}

func (r *HTMLRenderer) Render(data *parser.TableData) (string, error) {
	var html strings.Builder

	// Add style
	html.WriteString(r.Style)

	// Start table
	html.WriteString("<table>\n")

	// Add header row
	html.WriteString("  <tr>\n")
	for _, header := range data.Headers {
		html.WriteString(fmt.Sprintf("    <th>%s</th>\n", header))
	}
	html.WriteString("  </tr>\n")

	// Add data rows
	for _, row := range data.Rows {
		html.WriteString("  <tr>\n")
		for _, header := range data.Headers {
			html.WriteString(fmt.Sprintf("    <td>%s</td>\n", row[header]))
		}
		html.WriteString("  </tr>\n")
	}

	// Close table
	html.WriteString("</table>")

	return html.String(), nil
}
