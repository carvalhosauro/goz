package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"goz/internal/domain"
	"goz/internal/tui/theme"
)

type ListOpts struct {
	Width, Height int
	Tasks         []domain.Task
	Selected      int
	Adding        bool
	AddInput      *textinput.Model
}

// TaskList renders the column header, all task rows (with viewport-style
// clipping around the selected row), and an inline new-task input row.
func TaskList(t theme.Theme, opts ListOpts) string {
	if opts.Width <= 0 || opts.Height <= 0 {
		return ""
	}

	header := TaskRowHeader(t, opts.Width)
	bodyHeight := opts.Height - 1
	if opts.Adding {
		bodyHeight--
	}
	if bodyHeight < 1 {
		bodyHeight = 1
	}

	rows := make([]string, 0, len(opts.Tasks))
	for i, task := range opts.Tasks {
		rows = append(rows, TaskRow(t, task, i == opts.Selected, opts.Width))
	}

	// viewport-style clipping: keep selected row in view
	start := 0
	if len(rows) > bodyHeight {
		if opts.Selected >= bodyHeight {
			start = opts.Selected - bodyHeight + 1
		}
		if start+bodyHeight > len(rows) {
			start = len(rows) - bodyHeight
		}
	}
	end := start + bodyHeight
	if end > len(rows) {
		end = len(rows)
	}

	visible := rows[start:end]

	// pad short lists with empty filler
	for len(visible) < bodyHeight {
		visible = append(visible, lipgloss.NewStyle().Width(opts.Width).Render(""))
	}

	body := strings.Join(visible, "\n")
	parts := []string{header, body}

	if opts.Adding && opts.AddInput != nil {
		input := opts.AddInput.View()
		row := lipgloss.NewStyle().
			Foreground(t.Accent).
			Width(opts.Width).
			Render("  [ ] " + input)
		parts = append(parts, row)
	}

	return strings.Join(parts, "\n")
}
