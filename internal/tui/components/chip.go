package components

import (
	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

// Chip renders a small colored label. Charm theme uses dim brackets around
// the colored text since terminals can't fake CSS rounded pills.
func Chip(t theme.Theme, color lipgloss.Color, text string, faded bool) string {
	if text == "" {
		return ""
	}
	bracketStyle := lipgloss.NewStyle().Foreground(t.VeryDim)
	textStyle := lipgloss.NewStyle().Foreground(color)
	if faded {
		textStyle = textStyle.Faint(true)
		bracketStyle = bracketStyle.Faint(true)
	}
	return bracketStyle.Render("#") + textStyle.Render(text)
}
