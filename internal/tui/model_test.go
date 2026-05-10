package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"goz/internal/tui/components"
)

func sized(m Model, w, h int) Model {
	next, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return next.(Model)
}

func sendKey(m Model, key string) Model {
	var msg tea.KeyMsg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	case "tab":
		msg = tea.KeyMsg{Type: tea.KeyTab}
	case "space":
		msg = tea.KeyMsg{Type: tea.KeySpace}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	next, _ := m.Update(msg)
	return next.(Model)
}

func TestViewBelowMinSize(t *testing.T) {
	m := sized(NewModel(), 40, 10)
	v := m.View()
	if !strings.Contains(v, "needs at least") {
		t.Fatalf("expected size warning, got:\n%s", v)
	}
}

func TestViewRendersAllPanels(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	v := m.View()
	for _, want := range []string{
		"goz",
		"tasks",
		"today",
		"agent",
		"sidekick",
		"NORMAL",
		"Review Marcos",
	} {
		if !strings.Contains(v, want) {
			t.Errorf("expected %q in view, missing", want)
		}
	}
}

func TestToggleSelected(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	first := m.tasks[0]
	if first.Done {
		t.Fatal("seed first task should not be done")
	}
	m = sendKey(m, "x")
	if !m.tasks[0].Done {
		t.Fatal("expected first task done after toggle")
	}
}

func TestNavigation(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	if m.selected != 0 {
		t.Fatalf("selected=%d", m.selected)
	}
	m = sendKey(m, "j")
	if m.selected != 1 {
		t.Fatalf("after j: selected=%d", m.selected)
	}
	m = sendKey(m, "k")
	if m.selected != 0 {
		t.Fatalf("after k: selected=%d", m.selected)
	}
}

func TestDelete(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	before := len(m.tasks)
	m = sendKey(m, "d")
	if len(m.tasks) != before-1 {
		t.Fatalf("len before=%d after=%d", before, len(m.tasks))
	}
}

func TestAddTaskFlow(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	before := len(m.tasks)
	m = sendKey(m, "n")
	if !m.adding {
		t.Fatal("expected adding mode")
	}
	for _, ch := range "hello" {
		m = sendKey(m, string(ch))
	}
	m = sendKey(m, "enter")
	if m.adding {
		t.Fatal("expected adding cleared")
	}
	if len(m.tasks) != before+1 {
		t.Fatalf("len before=%d after=%d", before, len(m.tasks))
	}
	if m.tasks[len(m.tasks)-1].Text != "hello" {
		t.Fatalf("text=%s", m.tasks[len(m.tasks)-1].Text)
	}
}

func TestAddCancelEsc(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	before := len(m.tasks)
	m = sendKey(m, "n")
	m = sendKey(m, "esc")
	if m.adding {
		t.Fatal("adding should be cancelled")
	}
	if len(m.tasks) != before {
		t.Fatalf("len changed: before=%d after=%d", before, len(m.tasks))
	}
}

func TestTabFocus(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	if m.focus != components.ModeList {
		t.Fatal("default focus")
	}
	m = sendKey(m, "tab")
	if m.focus != components.ModeChat {
		t.Fatal("expected chat focus")
	}
	m = sendKey(m, "tab")
	if m.focus != components.ModeList {
		t.Fatal("expected list focus")
	}
}

func TestHelpToggle(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	m = sendKey(m, "?")
	if !m.showHelp {
		t.Fatal("expected help visible")
	}
	v := m.View()
	if !strings.Contains(v, "keybindings") {
		t.Fatal("expected help content")
	}
	m = sendKey(m, "?")
	if m.showHelp {
		t.Fatal("expected help hidden")
	}
}

func TestMoveTask(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	first := m.tasks[0].ID
	m = sendKey(m, "J")
	if m.tasks[1].ID != first {
		t.Fatalf("expected moved down, got order")
	}
	if m.selected != 1 {
		t.Fatalf("selected=%d", m.selected)
	}
}

func TestQuitKey(t *testing.T) {
	m := sized(NewModel(), 140, 40)
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected quit cmd")
	}
}
