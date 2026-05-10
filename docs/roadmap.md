# Roadmap

Phase-by-phase plan. Each phase teaches one concept and ships a working slice of the agent. Don't pre-implement future phases — that defeats the learning point.

**Status legend:** `[ ]` todo · `[~]` in progress · `[x]` done · `[-]` skipped/deferred

For locked decisions (charm only, EN, manual sort, etc.) see [decisions.md](./decisions.md).
For the target project layout see [design/layout.md](./design/layout.md).
For the stack reference see [architecture.md](./architecture.md).

---

## Phase overview

| # | Name | Goal | Status |
|---|------|------|--------|
| 0 | TUI foundation | Charm-themed TUI, in-memory CRUD | `[x]` |
| 1 | Persistence | SQLite store + migrations | `[ ]` |
| 2 | First AI call | Ollama client, prompt, streaming, temperature | `[ ]` |
| 3 | Tools / function calling | Tool registry, structured output, executor loop | `[ ]` |
| 4 | Agent loop & context | ReAct, conversation memory, summarization | `[ ]` |
| 5 | RAG | Embeddings, vector store, hybrid search | `[ ]` |
| 6 | Multi-agent orchestration | Planner / classifier / executor / critic graph | `[ ]` |
| 7 | Benchmark & eval | Golden dataset, metrics, A/B prompts | `[ ]` |
| 8 | Polish | MCP server, voice, plugins, sync | `[ ]` |

---

## Phase 0 — TUI foundation

**Scope:** charm-themed TUI, in-memory tasks (fixture seed), full keyboard CRUD, agent panel as visual placeholder. Zero persistence, zero AI. Output: a polished TUI binary you can run and play with.

### 0.1 Project scaffolding
- [x] `go mod init goz`
- [x] Directory layout (see [design/layout.md](./design/layout.md))
- [x] `Makefile` targets: `build`, `run`, `test`, `test-race`, `lint`, `tidy`, `fmt`, `vet`, `clean`
- [x] `.golangci.yml` strict config
- [x] `.gitignore`, `.editorconfig`
- [x] `README.md` skeleton (one-paragraph + run instructions)
- [x] `cmd/goz/main.go` minimal entrypoint that boots Bubble Tea
- [x] ADR `docs/adr/0001-stack-choices.md`

### 0.2 Theme system
- [x] `internal/tui/theme/theme.go` — `Theme` struct + `BoxChars`, `ChipKind`, `HeaderKind`
- [x] `internal/tui/theme/charm.go` — exact hex values from design
- [x] Lipgloss style helpers built from theme (`Style`, `Fore`, `DimStyle`, `AccentStyle`, `BorderColor`)

### 0.3 Domain model (in-memory)
- [x] `internal/domain/task.go` — `Task` struct (id, text, done, priority, estimate, due, tag, parentID, sortOrder, overdue)
- [x] `internal/domain/priority.go` — `P1|P2|P3|PNone` enum + parsing
- [x] `internal/domain/tag.go` — known tags
- [x] `internal/domain/seed.go` — fixture task list ported from `data.js` `TASKS`, translated to EN
- [x] Validation helpers + tests

### 0.4 In-memory store
- [x] `internal/store/store.go` — `Store` interface (`List`, `Get`, `Add`, `Toggle`, `Delete`, `Move`)
- [x] `internal/store/memory.go` — in-memory implementation seeded from fixture
- [x] Manual sort order: `Move(id, delta int)` repositions and renumbers `SortOrder`
- [x] Unit tests for memory store (Add/Toggle/Delete/Move/Get + clamps)

### 0.5 Components
- [x] `components/box.go` — bordered panel with `├─ title ─┤` cap + `titleRight`
- [x] `components/header.go` — pill header (`✦ goz` + path + datetime + status dots)
- [x] `components/statusbar.go` — mode label (NORMAL/CHAT) + dynamic keybinds + theme label
- [x] `components/chip.go` — colored `#tag` chip
- [x] `components/prio.go` — priority glyph (P1 danger / P2 warn / P3 dim / `··` for none)
- [x] `components/taskrow.go` — `[ ] P1 text  #tag  due  est` + selection marker + strikethrough on done
- [x] `components/tasklist.go` — header row + viewport-style clipping around the selected row + inline add row
- [x] `components/cursor.go` — block cursor (blink driven by parent)
- [x] `components/util.go` — `Truncate`, `PadLeft`, `PadRight`, `TagColor` helpers

### 0.6 Root model & key handling
- [x] `internal/tui/model.go` — root `Model` (store, theme, tasks, selected, focus, adding, addInput, chatInput, width, height, blink, now, msgs, showHelp)
- [x] `internal/tui/messages.go` — `tasksLoadedMsg`, `taskAddedMsg`, `taskMutatedMsg`, `errMsg`, `blinkMsg`, `clockMsg`
- [x] `internal/tui/keys.go` — `KeyMap` (Up, Down, Toggle, New, Delete, MoveUp, MoveDown, Tab, Quit, Help, Enter, Escape)
- [x] `Init` starts blink + clock tickers
- [x] `Update` routes messages and key events through focus-aware handlers
- [x] `View` composes header + body (2-col 1.45:1) + status bar with min-size guard
- [x] Inline new-task input (`n` triggers `textinput` row)
- [x] Resize handling (`tea.WindowSizeMsg` → recompute column widths)
- [x] Quit flow (`q` / `Ctrl+C`)

### 0.7 Agent panel placeholder
- [x] Right-side `Box` titled `agent · sidekick`, `titleRight=ollama · idle`
- [x] Render seed messages (welcome + 3 hints + muted "ai not wired" note, EN)
- [x] Chat input bar (visual; pressing Enter logs the user line + a muted "ai not wired up yet — phase 2." reply)
- [x] Quick-action buttons row (`/breakdown /estimate /prioritize /replan`) — clipped to fit width
- [x] `Tab` toggles focus list ↔ chat (border highlight switches)

### 0.8 Tests
- [x] Memory store unit tests (10 cases)
- [x] Domain validation tests + seed integrity test
- [x] Model behavior tests (View min-size, full render, navigate, toggle, add, delete, esc-cancel, Tab focus, help toggle, move, quit)
- [x] Snapshot dump via `goz --preview WxH` (saved at `internal/tui/testdata/snapshot_140x40.txt`)

### 0.9 Run & polish
- [x] `make run` launches binary
- [x] Min size guard (80x22) — prints `goz needs at least…` message centered
- [x] Help overlay (`?` toggles full keybind list)
- [x] `goz --preview` headless snapshot for visual verification
- [x] Verified charm visual matches mock (header pill, two boxes, sep lines, status bar, selection marker, tags, overdue glyph)

**Phase 0 done when:**
- Binary runs, renders charm theme matching mock
- All keys work: `j/k`, `x`, `n`, `d`, `Tab`, `q`, `?`, `J/K`
- New tasks survive within session (in-memory)
- Resize is smooth, no overflow
- All tests green, lint clean

---

## Phase 1 — Persistence (SQLite)

- [ ] Pick driver: `modernc.org/sqlite` (pure Go)
- [ ] `internal/store/migrations/0001_init.up.sql` (tasks table + parent_id self-ref + indexes + updated_at trigger)
- [ ] `internal/store/sqlite.go` implementing `Store`
- [ ] `golang-migrate` integration on boot
- [ ] XDG paths: DB at `$XDG_DATA_HOME/goz/goz.db`, config `$XDG_CONFIG_HOME/goz/config.yaml`
- [ ] `koanf` config loader + env overrides
- [ ] `slog` JSON handler to `$XDG_STATE_HOME/goz/goz.log`
- [ ] Backfill seed only if DB empty
- [ ] Integration tests with temp DB
- [ ] ADR 0002 (storage choice)

---

## Phase 2 — First AI call (Ollama)

- [ ] `internal/llm/ollama.go` HTTP client (`/api/generate` + streaming SSE)
- [ ] System prompt template + user prompt builder
- [ ] Wire chat input to Ollama, stream tokens to agent panel
- [ ] Temperature slider in TUI (live)
- [ ] Token counter (rough estimate)
- [ ] Cancellation via context (Esc cancels stream)
- [ ] Retry/backoff on transient errors
- [ ] Mock provider for tests

---

## Phase 3 — Tools & function calling

- [ ] `internal/agent/tool.go` — `Tool` interface (`Name`, `Schema`, `Execute`)
- [ ] JSON schema generation
- [ ] Tools: `create_task`, `update_status`, `delete_task`, `search_tasks`, `set_priority`, `set_estimate`, `set_due`
- [ ] Structured output via Ollama `format: "json"`
- [ ] Tool executor loop (LLM → parse → exec → feedback → LLM)
- [ ] Diff card UI: preview tool result, `A` apply / `R` reject / `E` edit
- [ ] Human-in-the-loop confirmation for destructive ops
- [ ] Dry-run mode

---

## Phase 4 — Agent loop & context management

- [ ] ReAct pattern (Thought → Action → Observation)
- [ ] Conversation memory tables (`conversations`, `messages`)
- [ ] Context window manager (truncate, keep system prompt)
- [ ] Hierarchical summarization when history > N tokens
- [ ] Max-step guard against runaway loops
- [ ] Trace viewer in TUI (step-by-step graph execution)

---

## Phase 5 — RAG over history

- [ ] `nomic-embed-text` via Ollama
- [ ] `sqlite-vec` extension (or alternative)
- [ ] Chunking strategy: task + description + comments
- [ ] Hybrid search: BM25 (FTS5) + vector → rerank
- [ ] `search_similar_tasks(query, k)` tool
- [ ] Background indexing pipeline (worker pool)

---

## Phase 6 — Multi-agent orchestration

- [ ] Custom graph engine: `Node`, `Edge`, `State`
- [ ] Agents: planner (T=0.7), classifier (T=0.1), executor (T=0.2), critic (T=0.3), summarizer (T=0.4)
- [ ] Routing: conditional edges based on state
- [ ] Visualization in TUI

---

## Phase 7 — Benchmark & eval

- [ ] Golden dataset (`bench/dataset.json`): inputs → expected tool calls
- [ ] Eval runner: `goz bench --model X --temp Y`
- [ ] Metrics: tool-selection accuracy, latency p50/p99, tokens/req
- [ ] A/B prompt comparison
- [ ] CSV/HTML report

---

## Phase 8 — Polish & extras

- [ ] MCP server (expose goz as MCP for Claude/Cursor)
- [ ] Voice input (whisper.cpp local)
- [ ] Recurring tasks via cron parsing through LLM ("every Monday 9am")
- [ ] Plugin system (Lua/Starlark tools)
- [ ] Multi-device sync (CRDTs)
