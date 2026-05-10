package components

import (
	"github.com/charmbracelet/lipgloss"

	"goz/internal/domain"
	"goz/internal/tui/theme"
)

// PrioGlyph renders a 2-cell priority indicator. None renders as "··".
func PrioGlyph(t theme.Theme, p domain.Priority) string {
	if p == domain.PNone {
		return lipgloss.NewStyle().Foreground(t.VeryDim).Render("··")
	}
	var c lipgloss.Color
	switch p {
	case domain.P1:
		c = t.Danger
	case domain.P2:
		c = t.Warn
	default:
		c = t.Dim
	}
	return lipgloss.NewStyle().Foreground(c).Bold(true).Render(string(p))
}
