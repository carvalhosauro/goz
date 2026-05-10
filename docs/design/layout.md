# Target project layout

The full layout `goz` is heading toward across all phases. Phase 0 ships everything outside `llm/`, `agent/`, `rag/`, `bench/`, `config/`, and the SQLite-only files in `store/`.

```
goz/
├── cmd/goz/main.go              # entry · --version · --preview
├── internal/
│   ├── tui/
│   │   ├── model.go             # Bubble Tea root (Model + Init/Update/View)
│   │   ├── messages.go          # tasksLoadedMsg, blinkMsg, clockMsg, ...
│   │   ├── keys.go              # KeyMap (Up, Down, Toggle, New, Delete, Tab, ...)
│   │   ├── agent.go             # right-panel placeholder (Phase 2 wires LLM)
│   │   ├── theme/
│   │   │   ├── theme.go         # Theme struct + BoxChars + ChipKind + HeaderKind
│   │   │   └── charm.go         # peach/mauve hex values
│   │   ├── components/
│   │   │   ├── box.go           # rounded panel + title cap + dashed sep
│   │   │   ├── header.go        # pill header (✦ goz · path · clock · status dots)
│   │   │   ├── statusbar.go     # NORMAL/CHAT block + dynamic keybinds
│   │   │   ├── chip.go          # #tag chip
│   │   │   ├── prio.go          # P1/P2/P3 / ··
│   │   │   ├── taskrow.go       # row renderer (checkbox · prio · text · tag · due · est)
│   │   │   ├── tasklist.go      # header row + viewport-style clipping + inline add row
│   │   │   ├── cursor.go        # block cursor
│   │   │   └── util.go          # Truncate, PadLeft, PadRight, TagColor
│   │   └── testdata/
│   │       └── snapshot_140x40.txt
│   ├── domain/
│   │   ├── task.go              # Task struct + Validate
│   │   ├── priority.go          # P1|P2|P3|PNone
│   │   ├── tag.go               # known tags
│   │   └── seed.go              # fixture data ported from the design mock
│   ├── store/
│   │   ├── store.go             # Store interface
│   │   ├── memory.go            # Phase 0 — in-memory
│   │   ├── sqlite.go            # Phase 1 — SQLite-backed
│   │   └── migrations/          # Phase 1 — golang-migrate file source
│   ├── llm/                     # Phase 2 — Ollama HTTP client + streaming
│   ├── agent/                   # Phase 3 — Tool registry + executor loop
│   ├── rag/                     # Phase 5 — chunking + embeddings + hybrid search
│   ├── bench/                   # Phase 7 — eval harness + golden dataset
│   ├── config/                  # Phase 1 — koanf-backed config + XDG paths
│   └── version/version.go
├── docs/
│   ├── README.md
│   ├── roadmap.md
│   ├── architecture.md
│   ├── decisions.md
│   ├── adr/
│   │   └── 0001-stack-choices.md
│   └── design/
│       ├── theme.md
│       └── layout.md
├── go.mod
├── Makefile
├── README.md
├── CLAUDE.md
├── .golangci.yml
├── .editorconfig
└── .gitignore
```
