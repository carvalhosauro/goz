package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

type BoxOpts struct {
	Width, Height int
	Title         string
	TitleRight    string
	Active        bool
}

// Box renders a rounded panel with a title row, a dashed separator, and a body.
// The title row mimics the design mock's `├─ title ─┤` cap inline (not embedded
// in the border) since lipgloss borders don't support title insertion natively.
func Box(t theme.Theme, opts BoxOpts, body string) string {
	if opts.Width < 4 || opts.Height < 4 {
		return ""
	}

	border := t.BorderColor(opts.Active)
	innerWidth := opts.Width - 2 // borders take 2 cells

	// title cap
	titleColor := t.Dim
	if opts.Active {
		titleColor = t.Accent
	}
	dim := lipgloss.NewStyle().Foreground(t.VeryDim)
	cap := dim.Render(t.Box.VL+t.Box.H+" ") +
		lipgloss.NewStyle().Foreground(titleColor).Render(opts.Title) +
		dim.Render(" "+t.Box.H+t.Box.VR)

	right := lipgloss.NewStyle().Foreground(t.Dim).Render(opts.TitleRight)

	// title row: " cap   ...   right "
	leftW := lipgloss.Width(cap)
	rightW := lipgloss.Width(right)
	pad := innerWidth - leftW - rightW - 2 // leading + trailing spaces
	if pad < 1 {
		pad = 1
	}
	titleRow := " " + cap + strings.Repeat(" ", pad) + right + " "
	titleRow = lipgloss.NewStyle().Width(innerWidth).Render(titleRow)

	// dashed separator row
	sep := strings.Repeat("┄", innerWidth)
	sepRow := lipgloss.NewStyle().Foreground(t.VeryDim).Render(sep)

	// body fills remaining height (height - 2 borders - title - sep)
	bodyHeight := opts.Height - 2 - 1 - 1
	if bodyHeight < 1 {
		bodyHeight = 1
	}
	bodyContent := lipgloss.NewStyle().Width(innerWidth).Height(bodyHeight).Render(body)

	inner := lipgloss.JoinVertical(lipgloss.Left, titleRow, sepRow, bodyContent)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Render(inner)
}
