package components

import (
	"github.com/charmbracelet/lipgloss"

	"goz/internal/domain"
	"goz/internal/tui/theme"
)

// Truncate cuts a string to fit within width cells, appending an ellipsis
// when content is dropped. Width is measured in display cells (lipgloss).
func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	r := []rune(s)
	for len(r) > 0 && lipgloss.Width(string(r))+1 > width {
		r = r[:len(r)-1]
	}
	return string(r) + "…"
}

// PadRight pads with spaces on the right to reach width cells.
func PadRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + spaces(width-w)
}

// PadLeft pads with spaces on the left to reach width cells.
func PadLeft(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return spaces(width-w) + s
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, n)
	for i := range out {
		out[i] = ' '
	}
	return string(out)
}

// TagColor maps a known tag to a theme color.
func TagColor(t theme.Theme, tag domain.Tag) lipgloss.Color {
	switch tag {
	case domain.TagEng:
		return t.Accent2
	case domain.TagDoc:
		return t.Info
	case domain.TagWork:
		return t.Accent
	case domain.TagPersonal:
		return t.Success
	case domain.TagPeople:
		return t.Warn
	case domain.TagLearn:
		return t.Info
	}
	return t.Dim
}
