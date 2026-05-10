package theme

import "github.com/charmbracelet/lipgloss"

type ChipKind int

const (
	ChipRounded ChipKind = iota
	ChipSquare
	ChipUnderline
)

type HeaderKind int

const (
	HeaderPill HeaderKind = iota
	HeaderTabs
	HeaderPlain
)

type BoxChars struct {
	TL, TR, BL, BR string
	H, V           string
	VL, VR         string
	Cross          string
}

type Theme struct {
	Name string

	Bg, Panel, Panel2 lipgloss.Color
	Fg, Dim, VeryDim  lipgloss.Color

	Accent, Accent2 lipgloss.Color
	Success         lipgloss.Color
	Danger          lipgloss.Color
	Warn            lipgloss.Color
	Info            lipgloss.Color

	Border       lipgloss.Color
	BorderActive lipgloss.Color
	RowSel       lipgloss.Color

	Box    BoxChars
	Chip   ChipKind
	Header HeaderKind
}

// Style returns a base style with theme foreground/background.
func (t Theme) Style() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Fg).Background(t.Bg)
}

func (t Theme) Fore(c lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c)
}

func (t Theme) DimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Dim)
}

func (t Theme) AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Accent)
}

func (t Theme) Accent2Style() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Accent2)
}

// Border returns the border color depending on active state.
func (t Theme) BorderColor(active bool) lipgloss.Color {
	if active {
		return t.BorderActive
	}
	return t.Border
}
