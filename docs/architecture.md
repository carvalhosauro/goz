# Architecture

How the layers fit, what to follow, what to avoid. Phase 0 reality below; later phases extend in clean lanes (see [roadmap.md](./roadmap.md)).

## Layering

```
cmd/goz/main.go        →  internal/tui  →  internal/store  →  internal/domain
                              └→  internal/tui/components  └→  internal/tui/theme
```

The split exists so later phases can swap implementations without rewriting the rest:

- **`internal/domain`** — pure types (`Task`, `Priority`, `Tag`, `Seed()`). No I/O, no UI. Phase 0 keeps `Estimate` and `Due` as free-form strings deliberately; later phases will type them.
- **`internal/store`** — `Store` interface (`List/Get/Add/Toggle/Delete/Move`) plus the `Memory` implementation. Phase 1 adds `sqlite.go` behind the same interface. **TUI code must depend on `store.Store`, never on a concrete type** — that is what makes the swap mechanical.
- **`internal/tui/theme`** — single `Theme` struct holding all colors, box characters, chip kind, header kind. `Charm()` is the only constructor today. **Components must accept a `theme.Theme` and never hardcode colors**, so the visual system stays swappable even though we only ship one theme.
- **`internal/tui/components`** — pure render functions (`Box`, `Header`, `StatusBar`, `TaskRow`, `TaskList`, `Chip`, `PrioGlyph`, `Cursor`) that take a theme + options and return a string. They never own state. `util.go` holds shared helpers (`Truncate`, `PadLeft/Right`, `TagColor`).
- **`internal/tui` (root)** — Bubble Tea MVU root: `Model`, `Update`, `View`, plus `keys.go` (`KeyMap`), `messages.go` (msg types), `agent.go` (right-panel placeholder).

## Bubble Tea conventions used here

- **Focus-aware key dispatch.** `handleKey` routes to `handleListKey`, `handleChatKey`, or `handleAddingKey` based on `Model.focus` and `Model.adding`. Adding a new key means picking the right handler — global keys (Tab, `?`, Quit) stay in `handleKey` itself; mode-specific keys go in their handler.
- **Self-rescheduling tickers.** `blinkTick` (530ms) and `clockTick` (1s) re-issue themselves from `Update`. Don't add new periodic state without following the same pattern — `tea.Tick` returns a one-shot Cmd.
- **Resize math lives in `View`.** `tea.WindowSizeMsg` only stores `width/height`; layout (1.45:1 split, body height = total − header − status) is computed at render time. Don't precompute layout in `Update`.
- **Min-size guard.** `View` short-circuits to a centered warning below 80×22. Any layout change must keep the guard intact.
- **Headless rendering.** `tui.RenderAt(w, h)` builds a Model, sends one `WindowSizeMsg`, and returns `View()`. Prefer this in tests and snapshots over driving a real `tea.Program`.

## Width-clipping discipline

Lipgloss's `.Width(w)` *wraps* when content exceeds `w`, which silently breaks fixed-height boxes. The right-panel learned this the hard way (input + actions wrapped one row each, pushing the agent box 3 rows past the list box).

Safe path:
- `padRow(s, w)` and `clip(s, w)` in `agent.go`
- `Truncate`, `PadLeft`, `PadRight` in `components/util.go`

**When rendering inside a fixed-size Box, ensure each row's visible width ≤ inner width before calling `.Width()`.** If you need both clip-on-overflow and pad-on-underflow, write it explicitly — don't hope `.Width` does what you want.

## Agent panel contract (`internal/tui/agent.go`)

Phase 0 is a visual placeholder: seed messages, a `textinput` chat bar, and quick-action buttons. Pressing Enter logs a muted "ai not wired — phase 2" reply.

The contract — which Phase 2 (Ollama) and Phase 3 (tools) hook into — is:
- Messages region (variable height, clipped from the bottom so latest is visible)
- Dashed separator (1 row)
- Input row (1 row, prompt + `textinput`)
- Quick-action buttons row (1 row, clipped to visible width)

Don't wire real LLM calls here. They go in `internal/llm/` (Phase 2) and `internal/agent/` (Phase 3).

## Headless preview

`./bin/goz --preview WxH` renders one frame to stdout without a TTY. Used for:

- snapshot files in `internal/tui/testdata/`
- showing visual changes in PR descriptions
- spot-checking layout from a non-TTY context (CI, agents, etc.)

Regenerate snapshots after intentional visual changes:

```sh
./bin/goz --preview 140x40 > internal/tui/testdata/snapshot_140x40.txt
```

## Stack reference

| Layer | Library | Notes |
|-------|---------|-------|
| TUI runtime | `bubbletea` | MVU |
| Styling | `lipgloss` | Theme composition; mind the `.Width()` wrapping pitfall above |
| Components | `bubbles` | `textinput`, `viewport`, `key`, `help` |
| DB | `modernc.org/sqlite` | Pure Go, no CGO (Phase 1) |
| Migrations | `golang-migrate/migrate/v4` | File source (Phase 1) |
| Config | `koanf` | YAML + env merge (Phase 1) |
| Logs | `log/slog` | stdlib, JSON handler |
| Tests | `testify` + golden | Snapshot rendering |
| LLM | Ollama HTTP API | Phase 2 |
| Embeddings | `nomic-embed-text` via Ollama | Phase 5 |
| Vector store | SQLite + `sqlite-vec` | Co-located with task data, hybrid w/ FTS5 |
| Build | `goreleaser` | Phase 8 |
