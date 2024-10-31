package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gowtham2003/gotable/pkg/parser"
	"github.com/gowtham2003/gotable/pkg/renderer"
)

type state int

const (
	stateInputFile state = iota
	stateInputFormat
	stateOutputFormat
	stateOutputFile
	stateStyle
	stateFormatOptions
	statePreview
	stateProcessing
	stateDone
)

type model struct {
	state            state
	filepicker       filepicker.Model
	formatList       list.Model
	styleList        list.Model
	outputInput      textinput.Model
	spinner          spinner.Model
	viewport         viewport.Model
	inputFile        string
	inputFormat      string
	outputFormat     string
	outputFile       string
	style            string
	preview          string
	err              error
	width            int
	height           int
	progress         progress.Model
	capabilities     FormatCapabilities
	colorEnabled     bool
	fontFamily       string
	tableWidth       int
	showStyleOptions bool
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF75B7")).
			MarginLeft(2)

	inputFormats = []list.Item{
		item{title: "JSON", desc: "JavaScript Object Notation"},
		item{title: "CSV", desc: "Comma Separated Values"},
		item{title: "Excel", desc: "Microsoft Excel Spreadsheet"},
		item{title: "HTML", desc: "HTML Table Format"},
	}

	outputFormats = []list.Item{
		item{title: "ASCII", desc: "ASCII Table Format"},
		item{title: "HTML", desc: "HTML Table Format"},
		item{title: "Excel", desc: "Microsoft Excel Spreadsheet"},
		item{title: "CSV", desc: "Comma Separated Values"},
		item{title: "JSON", desc: "JavaScript Object Notation"},
		item{title: "Markdown", desc: "Markdown Table Format"},
		item{title: "PNG", desc: "PNG Image Format"},
	}

	styleOptions = []list.Item{
		item{title: "Single", desc: "Single line borders"},
		item{title: "Double", desc: "Double line borders"},
		item{title: "Rounded", desc: "Rounded corners"},
		item{title: "Minimal", desc: "Minimal borders"},
	}
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Add new message types
type conversionFinishedMsg struct {
	output string
	err    error
}

type processingMsg float64

// Add format capabilities
type FormatCapabilities struct {
	SupportsStyle   bool
	SupportsColors  bool
	SupportsFonts   bool
	SupportsWidth   bool
	SupportsPreview bool
	FileExtension   string
}

var formatCapabilities = map[string]FormatCapabilities{
	"ASCII": {
		SupportsStyle:   true,
		SupportsPreview: true,
		FileExtension:   "txt",
	},
	"HTML": {
		SupportsStyle:   true,
		SupportsColors:  true,
		SupportsFonts:   true,
		SupportsWidth:   true,
		SupportsPreview: true,
		FileExtension:   "html",
	},
	"Markdown": {
		SupportsStyle:   true,
		SupportsPreview: true,
		FileExtension:   "md",
	},
	"CSV": {
		SupportsPreview: true,
		FileExtension:   "csv",
	},
	"JSON": {
		SupportsPreview: true,
		FileExtension:   "json",
	},
	"Excel": {
		SupportsColors: true,
		SupportsFonts:  true,
		SupportsWidth:  true,
		FileExtension:  "xlsx",
	},
	"PNG": {
		SupportsStyle:  true,
		SupportsColors: true,
		SupportsFonts:  true,
		SupportsWidth:  true,
		FileExtension:  "png",
	},
}

func (m model) processConversion() tea.Cmd {
	return func() tea.Msg {
		// Read input file
		input, err := os.ReadFile(m.inputFile)
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error reading input file: %v", err)}
		}

		// Create parser
		p, err := parser.NewParser(strings.ToLower(m.inputFormat))
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error creating parser: %v", err)}
		}

		// Parse input
		data, err := p.Parse(input)
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error parsing input: %v", err)}
		}

		// Create renderer with appropriate options
		r, err := renderer.NewRenderer(strings.ToLower(m.outputFormat))
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error creating renderer: %v", err)}
		}

		// Apply style options based on capabilities
		if styler, ok := r.(renderer.Styleable); ok && m.capabilities.SupportsStyle {
			styleOpts := renderer.StyleOptions{
				BorderStyle: strings.ToLower(m.style),
			}

			if m.capabilities.SupportsColors {
				styleOpts.ColorEnabled = m.colorEnabled
			}

			if m.capabilities.SupportsFonts {
				styleOpts.FontFamily = m.fontFamily
			}

			if m.capabilities.SupportsWidth {
				styleOpts.TableWidth = m.tableWidth
			}

			styler.SetStyle(styleOpts)
		}

		// Render output
		output, err := r.Render(data)
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error rendering output: %v", err)}
		}

		return conversionFinishedMsg{output: output}
	}
}

func (m model) saveOutput(output string) tea.Cmd {
	return func() tea.Msg {
		err := os.WriteFile(m.outputFile, []byte(output), 0644)
		if err != nil {
			return conversionFinishedMsg{err: fmt.Errorf("error saving output: %v", err)}
		}
		return nil
	}
}

func initialModel() model {
	m := model{
		state:       stateInputFile,
		filepicker:  filepicker.New(),
		formatList:  list.New(inputFormats, list.NewDefaultDelegate(), 0, 0),
		styleList:   list.New(styleOptions, list.NewDefaultDelegate(), 0, 0),
		outputInput: textinput.New(),
		spinner:     spinner.New(),
		progress:    progress.New(progress.WithDefaultGradient()),
	}

	m.filepicker.CurrentDirectory, _ = os.Getwd()
	m.formatList.Title = "Select Input Format"
	m.styleList.Title = "Select Style"
	m.outputInput.Placeholder = "Enter output file path..."
	m.outputInput.Focus()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.filepicker.Init(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// m.filepicker.SetSize(msg.Width, msg.Height-4)
		m.filepicker.Height = msg.Height - 4
		m.formatList.SetSize(msg.Width, msg.Height-4)
		m.styleList.SetSize(msg.Width, msg.Height-4)
		m.viewport = viewport.New(msg.Width, msg.Height-4)
		m.progress.Width = msg.Width - 4

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == stateStyle && m.capabilities.SupportsColors {
				m.colorEnabled = !m.colorEnabled
			}
		}

	case processingMsg:
		var progressCmd tea.Cmd
		m.progress.SetPercent(float64(msg))
		return m, progressCmd

	case conversionFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateDone
			return m, nil
		}
		m.preview = msg.output
		m.state = statePreview
		m.viewport.SetContent(msg.output)
		return m, nil

	case error:
		m.err = msg
		return m, nil
	}

	switch m.state {
	case stateInputFile:
		m.filepicker, cmd = m.filepicker.Update(msg)
		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.inputFile = path
			m.state = stateInputFormat
		}
		return m, cmd

	case stateInputFormat:
		m.formatList, cmd = m.formatList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				m.inputFormat = m.formatList.SelectedItem().(item).Title()
				m.formatList.SetItems(outputFormats)
				m.formatList.Title = "Select Output Format"
				m.state = stateOutputFormat
			}
		}
		return m, cmd

	case stateOutputFormat:
		m.formatList, cmd = m.formatList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				m.outputFormat = m.formatList.SelectedItem().(item).Title()
				m.capabilities = formatCapabilities[m.outputFormat]
				m.showStyleOptions = m.capabilities.SupportsStyle
				m.state = stateOutputFile
			}
		}
		return m, cmd

	case stateOutputFile:
		m.outputInput, cmd = m.outputInput.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				m.outputFile = m.outputInput.Value()
				m.state = stateStyle
			}
		}
		return m, cmd

	case stateStyle:
		if !m.showStyleOptions {
			m.state = stateProcessing
			return m, m.processConversion()
		}

		m.styleList, cmd = m.styleList.Update(msg)
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "enter" {
				m.style = m.styleList.SelectedItem().(item).Title()
				if m.capabilities.SupportsColors || m.capabilities.SupportsFonts {
					m.state = stateFormatOptions
				} else {
					m.state = stateProcessing
					return m, m.processConversion()
				}
			}
		}
		return m, cmd

	case stateFormatOptions:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "enter":
				m.state = stateProcessing
				return m, m.processConversion()
			case "tab":
				// Toggle between different options
				return m, nil
			}
		}

	case stateProcessing:
		if m.preview == "" {
			m.state = stateDone
			return m, m.processConversion()
		} else {
			m.state = stateDone
			return m, m.saveOutput(m.preview)
		}

	case statePreview:
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "y", "Y":
				m.state = stateProcessing
				return m, m.saveOutput(m.preview)
			case "n", "N":
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s strings.Builder

	switch m.state {
	case stateInputFile:
		s.WriteString(titleStyle.Render("Select Input File"))
		s.WriteString("\n\n")
		s.WriteString(m.filepicker.View())

	case stateInputFormat, stateOutputFormat:
		s.WriteString(m.formatList.View())

	case stateOutputFile:
		s.WriteString(titleStyle.Render("Enter Output File Path"))
		s.WriteString("\n\n")
		s.WriteString(m.outputInput.View())

	case stateStyle:
		if !m.showStyleOptions {
			m.state = stateProcessing
			return m.View()
		}
		s.WriteString(titleStyle.Render("Select Style"))
		s.WriteString("\n\n")
		s.WriteString(m.styleList.View())

	case stateFormatOptions:
		s.WriteString(titleStyle.Render("Format Options"))
		s.WriteString("\n\n")

		if m.capabilities.SupportsColors {
			s.WriteString(fmt.Sprintf("Colors enabled: %v [TAB to toggle]\n", m.colorEnabled))
		}

		if m.capabilities.SupportsFonts {
			s.WriteString(fmt.Sprintf("\nFont family: %s\n", m.fontFamily))
		}

		if m.capabilities.SupportsWidth {
			s.WriteString(fmt.Sprintf("\nTable width: %d\n", m.tableWidth))
		}

		s.WriteString("\nPress ENTER to continue")

	case stateProcessing:
		s.WriteString(fmt.Sprintf("\n\n  %s Converting...", m.spinner.View()))
		s.WriteString("\n\n")
		s.WriteString(m.progress.View())

	case statePreview:
		s.WriteString(titleStyle.Render("Preview"))
		s.WriteString("\n\n")
		s.WriteString(m.viewport.View())
		s.WriteString("\n\nSave this output? (y/n)")

	case stateDone:
		if m.err != nil {
			s.WriteString(errorStyle.Render(fmt.Sprintf("\n\nError: %v", m.err)))
			s.WriteString("\n\nPress q to quit")
		} else {
			s.WriteString(successStyle.Render("\n\nConversion completed successfully!"))
			s.WriteString(fmt.Sprintf("\n\nOutput saved to: %s", m.outputFile))
			s.WriteString("\n\nPress q to quit")
		}
	}

	// Add help text based on current state and capabilities
	s.WriteString(m.helpView())

	return s.String()
}

func (m model) helpView() string {
	var help strings.Builder
	help.WriteString("\n\n")

	switch m.state {
	case stateStyle:
		if m.showStyleOptions {
			help.WriteString("↑/↓: Select style • ENTER: Confirm")
		}
	case stateFormatOptions:
		help.WriteString("TAB: Toggle options • ENTER: Continue")
	case statePreview:
		if m.capabilities.SupportsPreview {
			help.WriteString("Y: Save • N: Cancel")
		}
	}

	help.WriteString(" • CTRL+C: Quit")
	return helpStyle.Render(help.String())
}

// Add success style
var successStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#04B575")).
	Bold(true).
	MarginTop(1)

// Add style for help text
var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#626262")).
	MarginTop(1)

func StartTUI() error {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}
	return nil
}
