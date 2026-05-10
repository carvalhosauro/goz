package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"goz/internal/domain"
	"goz/internal/store"
	"goz/internal/tui/components"
	"goz/internal/tui/theme"
)

const (
	minWidth  = 80
	minHeight = 22
)

type Model struct {
	store store.Store
	theme theme.Theme
	keys  KeyMap

	tasks    []domain.Task
	selected int

	focus  components.Mode
	adding bool

	addInput  textinput.Model
	chatInput textinput.Model

	msgs []agentMsg

	width, height int
	now           time.Time
	blink         bool
	showHelp      bool
	flash         string
}

// RenderAt builds a model, sizes it to w x h, and returns the rendered View.
// Useful for headless snapshots and visual checks.
func RenderAt(w, h int) string {
	m := NewModel()
	next, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return next.(Model).View()
}

func NewModel() Model {
	t := theme.Charm()

	add := textinput.New()
	add.Placeholder = "new task…"
	add.Prompt = ""
	add.CharLimit = 200

	chat := textinput.New()
	chat.Placeholder = "ask the agent · /breakdown /estimate /prioritize /replan"
	chat.Prompt = ""
	chat.CharLimit = 400

	s := store.NewMemory(domain.Seed())
	return Model{
		store:     s,
		theme:     t,
		keys:      DefaultKeys(),
		tasks:     s.List(),
		focus:     components.ModeList,
		addInput:  add,
		chatInput: chat,
		msgs:      seedAgentMessages(),
		now:       time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(blinkTick(), clockTick())
}

func blinkTick() tea.Cmd {
	return tea.Tick(530*time.Millisecond, func(t time.Time) tea.Msg { return blinkMsg(t) })
}

func clockTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return clockMsg(t) })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m2 := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = m2.Width
		m.height = m2.Height
		return m, nil

	case blinkMsg:
		m.blink = !m.blink
		return m, blinkTick()

	case clockMsg:
		m.now = time.Time(m2)
		return m, clockTick()

	case tea.KeyMsg:
		next, cmd := m.handleKey(m2)
		return next, cmd
	}

	// forward to focused input
	if m.adding {
		var cmd tea.Cmd
		m.addInput, cmd = m.addInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.focus == components.ModeChat && !m.adding {
		var cmd tea.Cmd
		m.chatInput, cmd = m.chatInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// global keys (always active)
	switch {
	case key.Matches(msg, m.keys.Quit) && !m.adding && m.focus != components.ModeChat:
		return m, tea.Quit
	case msg.String() == "ctrl+c":
		return m, tea.Quit
	}

	if m.showHelp {
		switch msg.String() {
		case "?", "esc", "q":
			m.showHelp = false
		}
		return m, nil
	}

	if m.adding {
		return m.handleAddingKey(msg)
	}

	if key.Matches(msg, m.keys.Tab) {
		if m.focus == components.ModeList {
			m.focus = components.ModeChat
			m.chatInput.Focus()
		} else {
			m.focus = components.ModeList
			m.chatInput.Blur()
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Help) {
		m.showHelp = true
		return m, nil
	}

	if m.focus == components.ModeChat {
		return m.handleChatKey(msg)
	}
	return m.handleListKey(msg)
}

func (m Model) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Down):
		if m.selected < len(m.tasks)-1 {
			m.selected++
		}
	case key.Matches(msg, m.keys.Up):
		if m.selected > 0 {
			m.selected--
		}
	case key.Matches(msg, m.keys.Toggle):
		if t := m.currentTask(); t != nil {
			_ = m.store.Toggle(t.ID)
			m.tasks = m.store.List()
		}
	case key.Matches(msg, m.keys.New):
		m.adding = true
		m.addInput.SetValue("")
		m.addInput.Focus()
	case key.Matches(msg, m.keys.Delete):
		if t := m.currentTask(); t != nil {
			_ = m.store.Delete(t.ID)
			m.tasks = m.store.List()
			if m.selected >= len(m.tasks) && m.selected > 0 {
				m.selected = len(m.tasks) - 1
			}
		}
	case key.Matches(msg, m.keys.MoveDown):
		if t := m.currentTask(); t != nil {
			_ = m.store.Move(t.ID, 1)
			m.tasks = m.store.List()
			if m.selected < len(m.tasks)-1 {
				m.selected++
			}
		}
	case key.Matches(msg, m.keys.MoveUp):
		if t := m.currentTask(); t != nil {
			_ = m.store.Move(t.ID, -1)
			m.tasks = m.store.List()
			if m.selected > 0 {
				m.selected--
			}
		}
	}
	return m, nil
}

func (m Model) handleAddingKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		text := strings.TrimSpace(m.addInput.Value())
		if text != "" {
			if _, err := m.store.Add(text); err == nil {
				m.tasks = m.store.List()
				m.selected = len(m.tasks) - 1
			}
		}
		m.adding = false
		m.addInput.Blur()
		return m, nil
	case "esc":
		m.adding = false
		m.addInput.Blur()
		return m, nil
	}
	var cmd tea.Cmd
	m.addInput, cmd = m.addInput.Update(msg)
	return m, cmd
}

func (m Model) handleChatKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		text := strings.TrimSpace(m.chatInput.Value())
		if text != "" {
			m.msgs = append(m.msgs, agentMsg{role: "user", text: text})
			m.msgs = append(m.msgs, agentMsg{role: "agent", muted: true, text: "(ai not wired up yet — phase 2.)"})
			m.chatInput.SetValue("")
		}
		return m, nil
	case "esc":
		m.focus = components.ModeList
		m.chatInput.Blur()
		return m, nil
	}
	var cmd tea.Cmd
	m.chatInput, cmd = m.chatInput.Update(msg)
	return m, cmd
}

func (m Model) currentTask() *domain.Task {
	if len(m.tasks) == 0 || m.selected < 0 || m.selected >= len(m.tasks) {
		return nil
	}
	t := m.tasks[m.selected]
	return &t
}

func (m Model) stats() components.Stats {
	var s components.Stats
	s.Total = len(m.tasks)
	for _, t := range m.tasks {
		switch {
		case t.Done:
			s.Done++
		default:
			s.Open++
			if t.Overdue {
				s.Overdue++
			}
		}
	}
	return s
}

func currentPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" && strings.HasPrefix(wd, home) {
		return "~" + strings.TrimPrefix(wd, home)
	}
	return wd
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}
	if m.width < minWidth || m.height < minHeight {
		msg := fmt.Sprintf("goz needs at least %dx%d (current %dx%d)", minWidth, minHeight, m.width, m.height)
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(m.theme.Danger).Render(msg))
	}

	if m.showHelp {
		return m.renderHelp()
	}

	headerH, statusH := 1, 1
	bodyH := m.height - headerH - statusH
	gap := 1
	leftW := (m.width - gap) * 145 / 245
	rightW := m.width - gap - leftW

	header := components.Header(m.theme, m.width, currentPath(), m.stats(), m.now)

	listInner := components.TaskList(m.theme, components.ListOpts{
		Width:    leftW - 2,
		Height:   bodyH - 4,
		Tasks:    m.tasks,
		Selected: m.selected,
		Adding:   m.adding,
		AddInput: &m.addInput,
	})
	st := m.stats()
	listBox := components.Box(m.theme, components.BoxOpts{
		Width:      leftW,
		Height:     bodyH,
		Title:      "tasks · today",
		TitleRight: fmt.Sprintf("%d/%d · %d overdue", st.Done, st.Total, st.Overdue),
		Active:     m.focus == components.ModeList,
	}, listInner)

	agentInner := renderAgent(
		m.theme,
		rightW-2,
		bodyH-4,
		m.msgs,
		m.chatInput,
		m.blink,
		m.focus == components.ModeChat,
	)
	agentBox := components.Box(m.theme, components.BoxOpts{
		Width:      rightW,
		Height:     bodyH,
		Title:      "agent · sidekick",
		TitleRight: "ollama · idle",
		Active:     m.focus == components.ModeChat,
	}, agentInner)

	body := lipgloss.JoinHorizontal(lipgloss.Top, listBox, " ", agentBox)

	binds := m.statusBinds()
	status := components.StatusBar(m.theme, m.width, m.focus, binds, "goz · charm")

	return lipgloss.JoinVertical(lipgloss.Left, header, body, status)
}

func (m Model) statusBinds() []components.Keybind {
	if m.focus == components.ModeChat {
		return []components.Keybind{
			{Key: "↵", Label: "send"},
			{Key: "Tab", Label: "list"},
			{Key: "Esc", Label: "exit"},
			{Key: "?", Label: "help"},
		}
	}
	if m.adding {
		return []components.Keybind{
			{Key: "↵", Label: "save"},
			{Key: "Esc", Label: "cancel"},
		}
	}
	return []components.Keybind{
		{Key: "j/k", Label: "nav"},
		{Key: "x", Label: "toggle"},
		{Key: "n", Label: "new"},
		{Key: "d", Label: "del"},
		{Key: "J/K", Label: "move"},
		{Key: "Tab", Label: "chat"},
		{Key: "?", Label: "help"},
		{Key: "q", Label: "quit"},
	}
}

func (m Model) renderHelp() string {
	rows := [][2]string{
		{"j / ↓", "move selection down"},
		{"k / ↑", "move selection up"},
		{"x / space", "toggle done"},
		{"n", "new task"},
		{"d", "delete task"},
		{"J", "move task down"},
		{"K", "move task up"},
		{"Tab", "switch focus list ↔ chat"},
		{"Enter", "send message (chat) / save (new task)"},
		{"Esc", "cancel new / leave chat"},
		{"?", "toggle this help"},
		{"q / Ctrl+C", "quit"},
	}
	keyStyle := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true).Width(14)
	descStyle := lipgloss.NewStyle().Foreground(m.theme.Fg)
	var lines []string
	lines = append(lines, lipgloss.NewStyle().Foreground(m.theme.Accent2).Bold(true).Render("goz · keybindings"))
	lines = append(lines, lipgloss.NewStyle().Foreground(m.theme.Dim).Render("press ? or esc to close"))
	lines = append(lines, "")
	for _, r := range rows {
		lines = append(lines, keyStyle.Render(r[0])+"  "+descStyle.Render(r[1]))
	}
	body := strings.Join(lines, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderActive).
		Padding(1, 3).
		Render(body)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
