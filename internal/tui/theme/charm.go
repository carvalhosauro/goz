package theme

import "github.com/charmbracelet/lipgloss"

// Charm returns the charm theme: peach + mauve on warm dark background.
func Charm() Theme {
	return Theme{
		Name:         "charm",
		Bg:           lipgloss.Color("#16131c"),
		Panel:        lipgloss.Color("#1d1925"),
		Panel2:       lipgloss.Color("#241e2d"),
		Fg:           lipgloss.Color("#ece4d6"),
		Dim:          lipgloss.Color("#6a6175"),
		VeryDim:      lipgloss.Color("#3d3645"),
		Accent:       lipgloss.Color("#f5a97f"),
		Accent2:      lipgloss.Color("#c6a0f6"),
		Success:      lipgloss.Color("#a6da95"),
		Danger:       lipgloss.Color("#ed8796"),
		Warn:         lipgloss.Color("#eed49f"),
		Info:         lipgloss.Color("#8aadf4"),
		Border:       lipgloss.Color("#3a3142"),
		BorderActive: lipgloss.Color("#f5a97f"),
		RowSel:       lipgloss.Color("#2c2434"),
		Box: BoxChars{
			TL: "╭", TR: "╮", BL: "╰", BR: "╯",
			H: "─", V: "│",
			VL: "├", VR: "┤",
			Cross: "┼",
		},
		Chip:   ChipRounded,
		Header: HeaderPill,
	}
}
