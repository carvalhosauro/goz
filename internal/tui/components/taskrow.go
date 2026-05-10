package components

import (
	"github.com/charmbracelet/lipgloss"

	"goz/internal/domain"
	"goz/internal/tui/theme"
)

// Column widths (excluding the flex text column).
const (
	colMarker = 2  // selection indicator
	colCheck  = 3  // [ ]
	colPrio   = 2  // P1 / ··
	colTag    = 12 // #engineering
	colDue    = 14 // ⚠ yesterday
	colEst    = 6  // 45m / —
	gap       = 1
)

// TaskRow renders a single task line at exactly `width` cells.
func TaskRow(t theme.Theme, task domain.Task, selected bool, width int) string {
	if width <= 0 {
		return ""
	}

	// selection marker
	marker := "  "
	if selected {
		marker = lipgloss.NewStyle().Foreground(t.Accent).Render("▍ ")
	}

	// checkbox
	checkbox := "[ ]"
	cbColor := t.Dim
	if task.Done {
		checkbox = "[x]"
		cbColor = t.Success
	}
	cb := lipgloss.NewStyle().Foreground(cbColor).Render(checkbox)

	// priority
	prio := PrioGlyph(t, task.Priority)

	// tag
	tag := ""
	if task.Tag != domain.TagNone {
		tag = Chip(t, TagColor(t, task.Tag), string(task.Tag), task.Done)
	}
	tagCell := PadLeft(tag, colTag)

	// due
	dueText := task.Due
	dueColor := t.Dim
	if task.Overdue && !task.Done {
		dueColor = t.Danger
		dueText = "⚠ " + dueText
	}
	dueCell := PadLeft(lipgloss.NewStyle().Foreground(dueColor).Render(dueText), colDue)

	// estimate
	est := task.Estimate
	if est == "" {
		est = "—"
	}
	estCell := PadLeft(lipgloss.NewStyle().Foreground(t.Dim).Render(est), colEst)

	// flex text width
	used := colMarker + colCheck + gap + colPrio + gap + colTag + gap + colDue + gap + colEst
	textWidth := width - used - gap*2 // leading + trailing pad
	if textWidth < 4 {
		textWidth = 4
	}

	textColor := t.Fg
	if task.Done {
		textColor = t.Dim
	}
	textStyle := lipgloss.NewStyle().Foreground(textColor)
	if task.Done {
		textStyle = textStyle.Strikethrough(true)
	}
	text := PadRight(textStyle.Render(Truncate(task.Text, textWidth)), textWidth)

	row := marker + cb + " " + prio + " " + text + " " + tagCell + " " + dueCell + " " + estCell

	if selected {
		row = lipgloss.NewStyle().Background(t.RowSel).Width(width).Render(row)
	} else {
		row = lipgloss.NewStyle().Width(width).Render(row)
	}
	return row
}

// TaskRowHeader renders a dim column header line.
func TaskRowHeader(t theme.Theme, width int) string {
	dim := lipgloss.NewStyle().Foreground(t.Dim)
	used := colMarker + colCheck + gap + colPrio + gap + colTag + gap + colDue + gap + colEst
	textWidth := width - used - gap*2
	if textWidth < 4 {
		textWidth = 4
	}
	row := "  " +
		dim.Render("   ") + " " +
		dim.Render("pr") + " " +
		dim.Render(PadRight("task", textWidth)) + " " +
		dim.Render(PadLeft("tag", colTag)) + " " +
		dim.Render(PadLeft("due", colDue)) + " " +
		dim.Render(PadLeft("est", colEst))
	return lipgloss.NewStyle().Width(width).Render(row)
}
