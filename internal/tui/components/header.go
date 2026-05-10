package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

type Stats struct {
	Total, Done, Open, Overdue int
}

// Header renders the pill-style top bar:
//
//	[✦ goz]  ~/path  ·  mon · 9 may · 14:23     ● 2 overdue  ● 5 done  ● 10 open
func Header(t theme.Theme, width int, path string, stats Stats, now time.Time) string {
	pill := lipgloss.NewStyle().
		Background(t.Accent).
		Foreground(t.Bg).
		Padding(0, 1).
		Bold(true).
		Render("✦ goz")

	pathStyled := lipgloss.NewStyle().Foreground(t.Fg).Render(path)
	dot := lipgloss.NewStyle().Foreground(t.VeryDim).Render("·")
	whenStyled := lipgloss.NewStyle().Foreground(t.Dim).
		Render(strings.ToLower(now.Format("Mon · 2 Jan · 15:04")))

	left := lipgloss.JoinHorizontal(lipgloss.Top,
		pill, " ", pathStyled, "  ", dot, "  ", whenStyled,
	)

	dotR := lipgloss.NewStyle().Foreground(t.Danger).Render("●")
	dotG := lipgloss.NewStyle().Foreground(t.Success).Render("●")
	dotA := lipgloss.NewStyle().Foreground(t.Accent).Render("●")
	dimText := lipgloss.NewStyle().Foreground(t.Dim).Render
	right := lipgloss.JoinHorizontal(lipgloss.Top,
		dotR, " ", dimText(fmt.Sprintf("%d overdue", stats.Overdue)),
		"   ",
		dotG, " ", dimText(fmt.Sprintf("%d done", stats.Done)),
		"   ",
		dotA, " ", dimText(fmt.Sprintf("%d open", stats.Open)),
	)

	leftW := lipgloss.Width(left)
	rightW := lipgloss.Width(right)
	spacer := width - leftW - rightW - 2 // 1-cell padding each side
	if spacer < 1 {
		spacer = 1
	}

	row := " " + left + strings.Repeat(" ", spacer) + right + " "
	return lipgloss.NewStyle().Width(width).Render(row)
}
