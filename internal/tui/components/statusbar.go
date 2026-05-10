package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

type Mode int

const (
	ModeList Mode = iota
	ModeChat
)

func (m Mode) Label() string {
	if m == ModeChat {
		return "CHAT"
	}
	return "NORMAL"
}

type Keybind struct {
	Key, Label string
}

func StatusBar(t theme.Theme, width int, mode Mode, binds []Keybind, themeLabel string) string {
	modeStyle := lipgloss.NewStyle().
		Foreground(t.Bg).
		Background(t.Accent).
		Bold(true).
		Padding(0, 1)
	modeBlock := modeStyle.Render(mode.Label())

	parts := []string{modeBlock}
	keyStyle := lipgloss.NewStyle().Foreground(t.Fg).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(t.Dim)
	for _, b := range binds {
		parts = append(parts, keyStyle.Render(b.Key)+" "+labelStyle.Render(b.Label))
	}

	left := strings.Join(parts, "  ")
	right := labelStyle.Render(themeLabel)

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	spacer := width - leftW - rightW - 2
	if spacer < 1 {
		spacer = 1
	}

	row := " " + left + strings.Repeat(" ", spacer) + right + " "
	return lipgloss.NewStyle().
		Background(t.Panel2).
		Width(width).
		Render(row)
}
