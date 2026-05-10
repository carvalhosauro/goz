# goz

> A terminal to-do list with a tiny AI agent that earns its keep.
> Built in Go. Written by hand. Learning out loud.

```
  ✦ goz  ~/repo/me/goz  ·  sat · 9 may · 01:40                          ● 1 overdue   ● 1 done   ● 9 open
╭────────────────────────────────────────────────────╮ ╭───────────────────────────────╮
│ ├─ tasks · today ─┤              1/10 · 1 overdue  │ │ ├─ agent · sidekick ─┤  …    │
│┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄│ │┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄│
│      pr task                    tag       due  est │ │ ◆ agent                       │
│▍ [ ] P1 Review Marcos's PR …    #eng    today  20m │ │ hey marina · 3 hints…         │
│  [ ] ·· Refactor billing …      #eng       —    —  │ │ …                             │
│  [ ] P1 Call accountant …  #personal ⚠ yesterday15m│ │                               │
│  …                                                 │ │┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄│
╰────────────────────────────────────────────────────╯ │ › ask the agent                │
                                                        ╰───────────────────────────────╯
  NORMAL  j/k nav  x toggle  n new  d del  J/K move  Tab chat  ? help  q quit       goz · charm
```

A learning rig disguised as a productivity app. The goal isn't to ship the next Things or Linear — it's to wire up an agent the slow, careful way: prompts first, then tools, then context, then RAG, then orchestration, then eval. Each phase lands a working binary. No mocks shipped to main. No big-bang.

## What it is right now

A charm-themed TUI todo list with manual reorder, an overdue indicator that judges you a little, and a "sidekick" panel where a real LLM will live in Phase 2.

- ✦ in-memory CRUD with manual sort (`Move(id, delta)` reorders and renumbers)
- ✦ priority colors, overdue glyph, strikethrough on done
- ✦ help overlay, min-size guard, headless `--preview` for snapshots
- ✦ 31 test cases (22 functions, 3 packages), race-clean, lint-clean
- ✦ zero AI yet — that's the point

## What it will be

`goz` is a phased project. Each phase teaches one concept and ships a working slice of the agent.

| Phase | Concept | Status |
|-------|---------|--------|
| 0 | TUI foundation (this) | ✅ |
| 1 | SQLite persistence + XDG paths | ⬜ next |
| 2 | First Ollama call · streaming · temperature | ⬜ |
| 3 | Tools & structured outputs | ⬜ |
| 4 | Agent loop · context · summarization | ⬜ |
| 5 | RAG over your task history | ⬜ |
| 6 | Multi-agent orchestration (custom graph) | ⬜ |
| 7 | Benchmark & eval harness | ⬜ |
| 8 | Polish — MCP server, sync, plugins | ⬜ |

Full checklists live at [docs/roadmap.md](./docs/roadmap.md).

## Run it

```sh
make run                       # boot the TUI
make build && ./bin/goz        # release-style build
./bin/goz --preview 140x40     # render a single frame to stdout
```

| Key | Action |
|-----|--------|
| `j` / `k` | move selection |
| `x` / `space` | toggle done |
| `n` | new task |
| `d` | delete |
| `J` / `K` | reorder task |
| `Tab` | switch focus list ↔ chat |
| `?` | help overlay |
| `q` / `Ctrl+C` | quit |

Min terminal size is 80×22. Below that, `goz` will tell you off.

## Development

```sh
make test                       # run all tests
make test-race                  # race detector
make lint                       # golangci-lint
make snapshot                   # regenerate TUI snapshot
make fmt                        # goimports + gofmt
make clean                      # nuke bin/ + coverage
```

The headless snapshot lives at `internal/tui/testdata/snapshot_140x40.txt`. Regenerate it with `make snapshot` after intentional visual changes.

## Stack

Go 1.25 · Bubble Tea + Lipgloss + Bubbles · pure-Go SQLite from Phase 1 · Ollama from Phase 2 · a custom multi-agent graph engine from Phase 6 — because rolling your own is the point. Reasoning lives in [ADR 0001](./docs/adr/0001-stack-choices.md).

## Docs

- [docs/architecture.md](./docs/architecture.md) — how the layers fit, the conventions to follow
- [docs/decisions.md](./docs/decisions.md) — what's locked and why
- [docs/roadmap.md](./docs/roadmap.md) — phase-by-phase tracking with checkboxes
- [docs/design/](./docs/design/) — theme reference and the target project layout
- [docs/adr/](./docs/adr/) — architecture decisions log
- [CLAUDE.md](./CLAUDE.md) — guidance for AI coding assistants working in this repo

## Why "goz"

It fell out of a design session as a placeholder and stuck. It sounds like a small reliable creature that does one thing well, which is the vibe.
