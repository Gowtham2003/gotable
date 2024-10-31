package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/gowtham2003/gotable/pkg/renderer"
	"github.com/gowtham2003/gotable/pkg/tui"
)

func main() {
	// Add flags
	cliMode := flag.Bool("cli", false, "Run in CLI mode")
	inputFormat := flag.String("if", "", "Input format (auto-detect by default)")
	outputFormat := flag.String("of", "", "Output format (auto-detect by default)")
	style := flag.String("style", "single", "Table style (single, double, rounded)")
	noHeader := flag.Bool("no-header", false, "Treat first row as data")
	help := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *cliMode {
		// Get input and output files from remaining arguments
		args := flag.Args()
		if len(args) < 2 {
			fmt.Println("Error: Input and output files are required")
			showHelp()
			os.Exit(1)
		}

		inputFile := args[0]
		outputFile := args[1]

		// Run CLI mode conversion
		if err := runCLIMode(cliOptions{
			inputFile:    inputFile,
			outputFile:   outputFile,
			inputFormat:  *inputFormat,
			outputFormat: *outputFormat,
			style:        *style,
			noHeader:     *noHeader,
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Run TUI mode by default
	if err := tui.StartTUI(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

type cliOptions struct {
	inputFile    string
	outputFile   string
	inputFormat  string
	outputFormat string
	style        string
	noHeader     bool
}

func runCLIMode(opts cliOptions) error {
	// Read input file
	input, err := os.ReadFile(opts.inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Auto-detect formats from file extensions if not specified
	if opts.inputFormat == "" {
		opts.inputFormat = detectFormat(opts.inputFile)
	}
	if opts.outputFormat == "" {
		opts.outputFormat = detectFormat(opts.outputFile)
	}

	// Create parser
	p, err := parser.NewParser(opts.inputFormat)
	if err != nil {
		return fmt.Errorf("failed to create parser: %v", err)
	}

	// Parse input
	data, err := p.Parse(input)
	if err != nil {
		return fmt.Errorf("failed to parse input: %v", err)
	}

	// Create renderer
	r, err := renderer.NewRenderer(opts.outputFormat)
	if err != nil {
		return fmt.Errorf("failed to create renderer: %v", err)
	}

	// Apply style if renderer supports it
	if styler, ok := r.(renderer.Styleable); ok {
		styler.SetStyle(renderer.StyleOptions{
			BorderStyle: opts.style,
			// NoHeader:    opts.noHeader,
		})
	}

	// Render output
	output, err := r.Render(data)
	if err != nil {
		return fmt.Errorf("failed to render output: %v", err)
	}

	// Write output file
	if err := os.WriteFile(opts.outputFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", opts.inputFile, opts.outputFile)
	return nil
}

func detectFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return "json"
	case ".csv":
		return "csv"
	case ".xlsx":
		return "excel"
	case ".html":
		return "html"
	case ".xml":
		return "xml"
	case ".md":
		return "markdown"
	case ".txt":
		return "ascii"
	case ".png":
		return "png"
	default:
		return "ascii" // default to ASCII for unknown formats
	}
}

func showHelp() {
	fmt.Println(`
GoTable - Universal Table Format Converter

Usage:
  gotable [flags] <input_file> <output_file>

Flags:
  -cli          Run in CLI mode
  -if string    Input format (auto-detect by default)
  -of string    Output format (auto-detect by default)
  -style string Table style (single, double, rounded) (default "single")
  -no-header    Treat first row as data
  -help         Show this help message

Supported Formats:
  Input:  json, csv, excel, html, xml
  Output: ascii, html, excel, csv, json, markdown, png

Examples:
  # Convert JSON to ASCII table
  gotable -cli input.json output.txt

  # Convert CSV to HTML with custom style
  gotable -cli -style rounded input.csv output.html

  # Convert Excel to Markdown without headers
  gotable -cli -no-header input.xlsx output.md

  # Convert with explicit formats
  gotable -cli -if json -of csv input.dat output.dat`)
}
