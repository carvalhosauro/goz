# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`goz` — a TUI to-do list with a local AI agent that manages tasks. Learning project for Go + agent orchestration (RAG, tools, context, eval). Phased roadmap in `docs/roadmap.md`; check it before starting work — it tracks per-phase scope and what is locked / done.

**Locked decisions (do not deviate):**
- Theme: **charm only** (peach `#f5a97f`, mauve `#c6a0f6`, warm dark bg). Tokyo and mono variants from the design mock are dropped.
- UI strings are **English**; seed content (names) may stay PT.
- **Manual sort order** for tasks (no auto re-rank). `Move(id, delta)` is the only reorder API.
- Module path is plain `goz` (no GitHub path yet).

**Phase status:** Phase 0 (TUI, in-memory) is complete. Phases 1+ are not implemented — do not pre-implement persistence, AI, or RAG until reaching that phase.

## Common commands

```sh
make build          # ldflags inject version, output bin/goz
make run            # go run ./cmd/goz
make test           # go test -count=1 ./...
make test-race      # go test -race -count=1 ./...
make lint           # golangci-lint run ./...
make fmt            # gofmt -s -w .
make vet            # go vet ./...
make tidy           # go mod tidy

go test ./internal/store -run TestMove          # single test
go test -v ./internal/tui -run TestAddTaskFlow  # single TUI test, verbose

./bin/goz --version            # print version
./bin/goz --preview 140x40     # render View headless to stdout (snapshot)
```

The headless preview is the primary way to verify visuals from a non-TTY context. The committed snapshot lives at `internal/tui/testdata/snapshot_140x40.txt`; regenerate it after intentional visual changes.

## Architecture (big picture)

The codebase follows a three-layer split that lets later phases swap implementations without touching the rest:

```
cmd/goz/main.go        →  internal/tui  →  internal/store  →  internal/domain
                              └→  internal/tui/components  └→  internal/tui/theme
```

- **`internal/domain`** — pure types (`Task`, `Priority`, `Tag`, `Seed()`). No I/O, no UI. Phase 0 keeps `Estimate` and `Due` as free-form strings deliberately; later phases will type them.
- **`internal/store`** — `Store` interface (`List/Get/Add/Toggle/Delete/Move`) plus the `Memory` implementation. Phase 1 will add `sqlite.go` behind the same interface; **TUI code must depend on `store.Store`, never on a concrete type**, so the swap is mechanical.
- **`internal/tui/theme`** — single `Theme` struct holding all colors, box characters, and chip/header style enums. `Charm()` is the only constructor today. **Components must accept a `theme.Theme` and never hardcode colors** — that is what keeps the visual system swappable.
- **`internal/tui/components`** — pure render functions (`Box`, `Header`, `StatusBar`, `TaskRow`, `TaskList`, `Chip`, `PrioGlyph`, `Cursor`) that take a `Theme` plus options and return a string. They never own state. `util.go` holds shared helpers (`Truncate`, `PadLeft/Right`, `TagColor`).
- **`internal/tui` (root)** — Bubble Tea MVU root: `Model`, `Update`, `View`, plus `keys.go` (`KeyMap`), `messages.go` (`blinkMsg`/`clockMsg`/etc.), and `agent.go` (right-panel placeholder).

### Bubble Tea conventions used here

- **Focus-aware key dispatch**: `handleKey` routes to `handleListKey`, `handleChatKey`, or `handleAddingKey` depending on `Model.focus` and `Model.adding`. Adding new keys means picking the right handler — global keys (Tab, `?`, Quit) live in `handleKey` itself, mode-specific keys live in their handler.
- **Tickers in `Init`**: `blinkTick` (530ms) and `clockTick` (1s) self-perpetuate by re-issuing themselves from `Update`. Don't add new periodic state without following the same self-rescheduling pattern.
- **Resize handling**: `tea.WindowSizeMsg` only updates `width/height`; layout math (1.45:1 split, body height = total − header − status) lives entirely in `View`. Don't precompute layout in `Update`.
- **Min-size guard**: `View` short-circuits to a centered warning below 80×22. Any layout change must keep that guard intact.
- **Headless rendering**: `tui.RenderAt(w, h)` builds a `Model`, sends one `WindowSizeMsg`, and returns `View()` output. Prefer this in tests and for snapshots over driving a real `tea.Program`.

### Agent panel (`agent.go`)

Phase 0 is a visual placeholder: seed messages, a `textinput` chat bar, and quick-action buttons. Pressing Enter logs a muted "ai not wired — phase 2" reply. **Do not wire real LLM calls here** — Phase 2 introduces `internal/llm/` and Phase 3 introduces `internal/agent/` with a tool registry. The placeholder UI shape (messages region + dashed separator + input row + actions row) is the contract those phases will hook into.

### Width clipping discipline

Lipgloss's `.Width(w)` *wraps* when content exceeds `w`, which silently breaks fixed-height boxes. The right-panel (`renderAgent`) found this the hard way. Helpers `padRow` and `clip` in `agent.go`, plus `Truncate`/`PadLeft`/`PadRight` in `components/util.go`, are the safe path. **When rendering inside a fixed-size Box, ensure each row's visible width ≤ inner width before calling `.Width()`.**

## Testing notes

- `internal/tui/model_test.go` drives the model directly via `tea.KeyMsg` and `tea.WindowSizeMsg`. Use the local `sized()` and `sendKey()` helpers when adding cases — they handle the `Model` type assertion.
- Visual tests are substring assertions on `View()` output, not pixel diffs. The committed snapshot is for human review and regeneration after intentional visual changes, not automated comparison (yet).
- Store tests live next to the store; TUI tests live in `internal/tui`. Keep them separated.

## Phase boundaries

Before starting work, check `docs/roadmap.md` for the current phase. Common foot-guns:
- Phase 1 adds SQLite behind `store.Store` and introduces XDG paths (`$XDG_DATA_HOME/goz/goz.db`, `$XDG_CONFIG_HOME/goz/config.yaml`, `$XDG_STATE_HOME/goz/goz.log`). Don't write files to `~/` directly.
- Phase 2 brings Ollama HTTP + streaming. Don't import an LLM SDK before then.
- Phase 5 brings RAG (`sqlite-vec` + `nomic-embed-text`). The vector store sits next to tasks in the same SQLite file.
- Phase 6 builds a custom graph engine (Node/Edge/State) for multi-agent orchestration. **Do not pull in `langchaingo` / `eino` / similar** — building it is a learning goal of the project (see ADR 0001).

## ADRs

`docs/adr/` holds architecture decisions. Read existing ADRs before changing stack choices, and add a new ADR (next number) when making a non-trivial decision (storage swap, framework change, agent topology, etc.).
