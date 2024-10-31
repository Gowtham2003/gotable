# GoTable - Universal Table Format Converter

![GoTable Demo](assets/demo.gif)

A powerful and flexible command-line tool for converting between various table formats with an interactive Terminal User Interface (TUI).

## Features

### Supported Input Formats

- JSON
- CSV
- Excel (XLSX)
- HTML
- XML

### Supported Output Formats

- ASCII Table
- HTML
- Excel (XLSX)
- CSV
- JSON
- Markdown
- PNG Image

### Key Features

- 🖼️ Interactive TUI with real-time preview
- 🎨 Customizable styling options
- 📊 Multiple border styles
- 🎯 Format-specific customization
- 🚀 Batch processing support
- 💾 Auto file extension handling
- 🎭 Light/Dark theme support

## Installation

```bash
# Using go install
go install github.com/gowtham2003/gotable@latest

# Or clone and build
git clone https://github.com/gowtham2003/gotable.git
cd gotable
go build

```

## Usage

### Interactive TUI Mode (Default)

```bash
gotable
```

### Command Line Mode

```bash
gotable -i input.json -o output.csv
```

### Examples

Convert JSON to ASCII Table:

```bash
gotable input.json output.txt
```

Convert Excel to Markdown:

```bash
gotable input.xlsx output.md
```

### Format-Specific Features

#### ASCII Table

- Single, double, or rounded borders
- Custom width
- Unicode support

#### HTML

- Customizable CSS styles
- Responsive design
- Color themes
- Custom fonts

#### Excel

- Auto-column width
- Header styling
- Custom fonts
- Cell formatting

#### PNG

- Custom dimensions
- Font selection
- Color schemes
- Border styles

## Configuration

### Style Options

Different output formats support various styling options:

| Format   | Borders | Colors | Fonts | Width | Preview |
| -------- | ------- | ------ | ----- | ----- | ------- |
| ASCII    | ✅      | ❌     | ❌    | ✅    | ✅      |
| HTML     | ✅      | ✅     | ✅    | ✅    | ✅      |
| Markdown | ✅      | ❌     | ❌    | ❌    | ✅      |
| CSV      | ❌      | ❌     | ❌    | ❌    | ✅      |
| JSON     | ❌      | ❌     | ❌    | ❌    | ✅      |
| Excel    | ✅      | ✅     | ✅    | ✅    | ❌      |
| PNG      | ✅      | ✅     | ✅    | ✅    | ✅      |

## Development

### Prerequisites

- Go 1.16 or higher
- Required dependencies:
  ```bash
  go get github.com/charmbracelet/bubbletea
  go get github.com/charmbracelet/bubbles
  go get github.com/charmbracelet/lipgloss
  go get github.com/xuri/excelize/v2
  ```

### Project Structure

```bash
gotable/
├── cmd/
│ └── gotable/
│ └── main.go
├── pkg/
│ ├── parser/
│ │ ├── parser.go
│ │ ├── json_parser.go
│ │ ├── csv_parser.go
│ │ └── ...
│ ├── renderer/
│ │ ├── renderer.go
│ │ ├── ascii_renderer.go
│ │ ├── html_renderer.go
│ │ └── ...
│ └── tui/
│ ├── tui.go
│ └── styles.go
├── test/
│ └── integration_test.go
├── go.mod
├── go.sum
└── README.md
```

### Running Tests

```bash

# Run all tests
go test ./...
# Run with coverage
go test ./... -cover
# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Excelize](https://github.com/xuri/excelize) - Excel file handling
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components

## Support

For support, please open an issue in the GitHub issue tracker or contact the maintainers.

---

Made with ❤️ by Gowtham
