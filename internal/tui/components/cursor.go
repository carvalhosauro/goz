package components

import (
	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

// Cursor renders a single-cell block cursor. Caller handles blink state.
func Cursor(t theme.Theme, on bool) string {
	if !on {
		return " "
	}
	return lipgloss.NewStyle().Background(t.Accent).Render(" ")
}
