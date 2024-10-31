package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/gowtham2003/gotable/pkg/renderer"
)

type InteractiveMode struct {
	reader *bufio.Reader
}

type InteractiveOptions struct {
	InputFile      string
	InputFormat    string
	OutputFile     string
	OutputFormat   string
	StyleOptions   renderer.StyleOptions
	PreviewEnabled bool
}

func New() *InteractiveMode {
	return &InteractiveMode{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (im *InteractiveMode) Run() error {
	options := &InteractiveOptions{}

	fmt.Println("=== Interactive Table Converter ===")

	// Get input file
	if err := im.getInputFile(options); err != nil {
		return err
	}

	// Get input format
	if err := im.getInputFormat(options); err != nil {
		return err
	}

	// Get output format
	if err := im.getOutputFormat(options); err != nil {
		return err
	}

	// Get style options
	if err := im.getStyleOptions(options); err != nil {
		return err
	}

	// Preview option
	if err := im.getPreviewOption(options); err != nil {
		return err
	}

	// Process the conversion
	return im.processConversion(options)
}

func (im *InteractiveMode) getInputFile(options *InteractiveOptions) error {
	fmt.Print("Enter input file path: ")
	input, err := im.reader.ReadString('\n')
	if err != nil {
		return err
	}

	options.InputFile = strings.TrimSpace(input)
	if _, err := os.Stat(options.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", options.InputFile)
	}
	return nil
}

func (im *InteractiveMode) getInputFormat(options *InteractiveOptions) error {
	fmt.Println("\nAvailable input formats:")
	fmt.Println("1. JSON")
	fmt.Println("2. CSV")
	fmt.Println("3. Excel")
	fmt.Println("4. HTML")
	fmt.Print("Select input format (1-4): ")

	input, err := im.reader.ReadString('\n')
	if err != nil {
		return err
	}

	switch strings.TrimSpace(input) {
	case "1":
		options.InputFormat = "json"
	case "2":
		options.InputFormat = "csv"
	case "3":
		options.InputFormat = "xlsx"
	case "4":
		options.InputFormat = "html"
	default:
		return fmt.Errorf("invalid input format selection")
	}
	return nil
}

func (im *InteractiveMode) getOutputFormat(options *InteractiveOptions) error {
	fmt.Println("\nAvailable output formats:")
	fmt.Println("1. ASCII Table")
	fmt.Println("2. HTML")
	fmt.Println("3. Excel")
	fmt.Println("4. CSV")
	fmt.Println("5. JSON")
	fmt.Println("6. Markdown")
	fmt.Println("7. PNG Image")
	fmt.Print("Select output format (1-7): ")

	input, err := im.reader.ReadString('\n')
	if err != nil {
		return err
	}

	switch strings.TrimSpace(input) {
	case "1":
		options.OutputFormat = "ascii"
	case "2":
		options.OutputFormat = "html"
	case "3":
		options.OutputFormat = "xlsx"
	case "4":
		options.OutputFormat = "csv"
	case "5":
		options.OutputFormat = "json"
	case "6":
		options.OutputFormat = "markdown"
	case "7":
		options.OutputFormat = "png"
	default:
		return fmt.Errorf("invalid output format selection")
	}

	fmt.Print("Enter output file path: ")
	input, err = im.reader.ReadString('\n')
	if err != nil {
		return err
	}
	options.OutputFile = strings.TrimSpace(input)
	return nil
}

func (im *InteractiveMode) getStyleOptions(options *InteractiveOptions) error {
	fmt.Println("\nStyle Options:")

	// Border style
	fmt.Println("Select border style:")
	fmt.Println("1. Single")
	fmt.Println("2. Double")
	fmt.Println("3. Rounded")
	fmt.Print("Choose style (1-3): ")

	input, err := im.reader.ReadString('\n')
	if err != nil {
		return err
	}

	switch strings.TrimSpace(input) {
	case "1":
		options.StyleOptions.BorderStyle = "single"
	case "2":
		options.StyleOptions.BorderStyle = "double"
	case "3":
		options.StyleOptions.BorderStyle = "rounded"
	default:
		options.StyleOptions.BorderStyle = "single"
	}

	// Color options (if supported by output format)
	if options.OutputFormat == "html" || options.OutputFormat == "png" {
		fmt.Print("Enable colored output? (y/n): ")
		input, err = im.reader.ReadString('\n')
		if err != nil {
			return err
		}
		options.StyleOptions.ColorEnabled = strings.ToLower(strings.TrimSpace(input)) == "y"
	}

	return nil
}

func (im *InteractiveMode) getPreviewOption(options *InteractiveOptions) error {
	fmt.Print("\nPreview output before saving? (y/n): ")
	input, err := im.reader.ReadString('\n')
	if err != nil {
		return err
	}
	options.PreviewEnabled = strings.ToLower(strings.TrimSpace(input)) == "y"
	return nil
}

func (im *InteractiveMode) processConversion(options *InteractiveOptions) error {
	// Read input file
	input, err := os.ReadFile(options.InputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	// Create parser
	p, err := parser.NewParser(options.InputFormat)
	if err != nil {
		return fmt.Errorf("error creating parser: %v", err)
	}

	// Parse input
	data, err := p.Parse(input)
	if err != nil {
		return fmt.Errorf("error parsing input: %v", err)
	}

	// Create renderer
	r, err := renderer.NewRenderer(options.OutputFormat)
	if err != nil {
		return fmt.Errorf("error creating renderer: %v", err)
	}

	// Apply style options if renderer supports it
	if styler, ok := r.(renderer.Styleable); ok {
		styler.SetStyle(options.StyleOptions)
	}

	// Render output
	output, err := r.Render(data)
	if err != nil {
		return fmt.Errorf("error rendering output: %v", err)
	}

	// Preview if enabled
	if options.PreviewEnabled {
		fmt.Println("\nPreview:")
		fmt.Println(output)
		fmt.Print("\nSave this output? (y/n): ")
		input, err := im.reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			return fmt.Errorf("operation cancelled by user")
		}
	}

	// Write to output file
	if err := os.WriteFile(options.OutputFile, []byte(output), 0644); err != nil {
		return fmt.Errorf("error writing output file: %v", err)
	}

	fmt.Printf("\nSuccessfully converted %s to %s\n", options.InputFile, options.OutputFile)
	return nil
}
