package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"goz/internal/tui"
	"goz/internal/version"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Println("goz", version.Version)
			return
		case "--preview":
			runPreview()
			return
		}
	}

	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "goz: %v\n", err)
		os.Exit(1)
	}
}

// runPreview renders the View at a fixed size and writes the result to stdout.
// Used for snapshotting and quick visual checks without a TTY.
//
// Usage:
//
//	goz --preview          # default 140x40
//	goz --preview 120x32
func runPreview() {
	w, h := 140, 40
	if len(os.Args) >= 3 {
		parts := strings.Split(os.Args[2], "x")
		if len(parts) == 2 {
			if iw, err := strconv.Atoi(parts[0]); err == nil {
				w = iw
			}
			if ih, err := strconv.Atoi(parts[1]); err == nil {
				h = ih
			}
		}
	}
	fmt.Print(tui.RenderAt(w, h))
}
