package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	special   = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}

	divider = lipgloss.NewStyle().
		Foreground(subtle).
		Padding(0, 1).
		Render("•")

	urlStyle = lipgloss.NewStyle().Foreground(special).Underline(true)

	activeButtonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFF")).
				Background(highlight).
				Padding(0, 3).
				MarginRight(2)

	inactiveButtonStyle = lipgloss.NewStyle().
				Foreground(subtle).
				Padding(0, 3).
				MarginRight(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			MarginTop(1)
)
