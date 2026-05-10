package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up, Down         key.Binding
	Toggle           key.Binding
	New              key.Binding
	Delete           key.Binding
	MoveUp, MoveDown key.Binding
	Tab              key.Binding
	Quit             key.Binding
	Help             key.Binding
	Enter, Escape    key.Binding
}

func DefaultKeys() KeyMap {
	return KeyMap{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k/↑", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j/↓", "down")),
		Toggle:   key.NewBinding(key.WithKeys("x", " "), key.WithHelp("x", "toggle")),
		New:      key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new")),
		Delete:   key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
		MoveUp:   key.NewBinding(key.WithKeys("K"), key.WithHelp("K", "move up")),
		MoveDown: key.NewBinding(key.WithKeys("J"), key.WithHelp("J", "move down")),
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "focus")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Enter:    key.NewBinding(key.WithKeys("enter")),
		Escape:   key.NewBinding(key.WithKeys("esc")),
	}
}
