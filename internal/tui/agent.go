package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"goz/internal/tui/theme"
)

type agentMsg struct {
	role  string // "agent" | "user"
	text  string
	muted bool
}

func seedAgentMessages() []agentMsg {
	return []agentMsg{
		{role: "agent", muted: true, text: "hey marina · 3 hints for today:"},
		{role: "agent", text: "1. \"Call accountant\" is overdue · /replan to move"},
		{role: "agent", text: "2. \"Refactor billing\" is large · /breakdown"},
		{role: "agent", text: "3. 3 tasks without estimate · /estimate"},
		{role: "agent", muted: true, text: "(ai not wired — phase 2 brings ollama.)"},
	}
}

type quickAction struct {
	label string
	cmd   string
}

var quickActions = []quickAction{
	{"breakdown", "break down a task"},
	{"estimate", "estimate my tasks"},
	{"prioritize", "prioritize my list"},
	{"replan", "replan overdue"},
}

// renderAgent builds the inner content for the agent box.
func renderAgent(t theme.Theme, width, height int, msgs []agentMsg, chat textinput.Model, _ bool, _ bool) string {
	if width < 8 || height < 4 {
		return ""
	}

	sep := lipgloss.NewStyle().Foreground(t.VeryDim).Render(strings.Repeat("┄", width))
	inputRow := renderAgentInput(t, width, chat)
	actionsRow := renderQuickActions(t, width)

	// Layout: msgs + sep + input + actions  ⇒  msgs height = total - 3
	msgArea := height - 3
	if msgArea < 1 {
		msgArea = 1
	}
	msgsBlock := renderAgentMessages(t, width, msgArea, msgs)

	return lipgloss.JoinVertical(lipgloss.Left, msgsBlock, sep, inputRow, actionsRow)
}

func renderAgentMessages(t theme.Theme, width, height int, msgs []agentMsg) string {
	lines := make([]string, 0, len(msgs)*2)
	for i, msg := range msgs {
		var head lipgloss.Style
		var label string
		if msg.role == "user" {
			head = lipgloss.NewStyle().Foreground(t.Accent)
			label = "› you"
		} else {
			head = lipgloss.NewStyle().Foreground(t.Accent2)
			label = "◆ agent"
		}
		lines = append(lines, padRow(head.Render(label), width))

		bodyColor := t.Fg
		if msg.muted {
			bodyColor = t.Dim
		}
		bodyStyle := lipgloss.NewStyle().Foreground(bodyColor)
		for _, raw := range strings.Split(msg.text, "\n") {
			lines = append(lines, padRow(bodyStyle.Render(clip(raw, width)), width))
		}
		if i < len(msgs)-1 {
			lines = append(lines, padRow("", width))
		}
	}

	if len(lines) > height {
		lines = lines[len(lines)-height:]
	}
	for len(lines) < height {
		lines = append(lines, padRow("", width))
	}
	return strings.Join(lines, "\n")
}

func renderAgentInput(t theme.Theme, width int, chat textinput.Model) string {
	prompt := lipgloss.NewStyle().Foreground(t.Accent).Render("› ")
	chat.Width = width - 3
	if chat.Width < 4 {
		chat.Width = 4
	}
	row := prompt + chat.View()
	return padRow(row, width)
}

func renderQuickActions(t theme.Theme, width int) string {
	style := lipgloss.NewStyle().
		Foreground(t.Dim).
		Background(t.Panel2).
		Padding(0, 1)
	var b strings.Builder
	used := 0
	for i, a := range quickActions {
		piece := style.Render("/" + a.label)
		w := lipgloss.Width(piece)
		if i > 0 {
			w++ // single-space gap
		}
		if used+w > width {
			break
		}
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(piece)
		used += w
	}
	return padRow(b.String(), width)
}

// padRow pads a string with trailing spaces (or clips) to exactly `width` cells.
func padRow(s string, width int) string {
	w := lipgloss.Width(s)
	if w == width {
		return s
	}
	if w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return clip(s, width)
}

// clip truncates by visible cells. Naive — assumes no ANSI escapes (caller's
// responsibility to clip pre-styling for known short strings).
func clip(s string, width int) string {
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
