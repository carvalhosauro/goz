# Charm theme reference

The single source of truth for `goz` colors and glyphs. Mirrors `internal/tui/theme/charm.go` — if these diverge, the Go file wins.

## Palette

| Role | Hex | Notes |
|------|-----|-------|
| `bg` | `#16131c` | Warm dark canvas |
| `panel` | `#1d1925` | Box body |
| `panel2` | `#241e2d` | Status bar / button background |
| `fg` | `#ece4d6` | Primary text |
| `dim` | `#6a6175` | Secondary text, separators |
| `veryDim` | `#3d3645` | Tertiary, decorative dashes |
| `accent` | `#f5a97f` | Peach — primary highlight, selection, active border |
| `accent2` | `#c6a0f6` | Mauve — secondary highlight, agent label |
| `success` | `#a6da95` | Done state |
| `danger` | `#ed8796` | P1 priority, overdue, warnings |
| `warn` | `#eed49f` | P2 priority |
| `info` | `#8aadf4` | Informational tags |
| `border` | `#3a3142` | Inactive box border |
| `borderActive` | `#f5a97f` | Focused box border (= accent) |
| `rowSel` | `#2c2434` | Selected task row background |

## Glyphs & box drawing

| Use | Char |
|-----|------|
| Box top-left / top-right | `╭` `╮` |
| Box bottom-left / bottom-right | `╰` `╯` |
| Box horizontal / vertical | `─` `│` |
| Title cap left / right | `├` `┤` |
| Dashed separator | `┄` |
| Selection marker | `▍ ` |
| Overdue indicator | `⚠ ` |
| Priority none | `··` |
| Logo | `✦ goz` |

## Component style notes

- **Header**: pill-style. Brand pill uses `accent` background, `bg` foreground, bold, 1-cell horizontal padding.
- **Box**: rounded border using the chars above. Title rendered inline as `├─ title ─┤` (lipgloss has no native title-on-border insertion).
- **Chip**: `#tag` — `veryDim` `#` glyph, tag-specific color text. Faded (50%) when the row is `done`.
- **Tag colors**: `eng → accent2`, `doc → info`, `work → accent`, `personal → success`, `people → warn`, `learn → info`.
- **Status bar**: `panel2` background. `NORMAL` / `CHAT` mode block uses `accent` background, `bg` foreground, bold.
- **Cursor**: 1-cell `accent` background, blinks at 530ms.

## Typography

JetBrains Mono in the design mock. The TUI itself inherits the user's terminal font; `lipgloss.Bold(true)` is used sparingly for emphasis.
