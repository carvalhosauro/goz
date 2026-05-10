# ADR 0001 — Stack choices

Date: 2026-05-09
Status: Accepted

## Context

goz is a learning project: a terminal TUI to-do list with a local AI agent. Goals: practice idiomatic Go, then explore agent orchestration, tools, RAG, and eval. Phase 0 is the TUI foundation; later phases add SQLite persistence, an Ollama-based agent, RAG, multi-agent orchestration, and a benchmark harness.

## Decision

| Concern | Choice | Reason |
|---------|--------|--------|
| TUI runtime | `bubbletea` | Mature MVU framework. De facto standard in Go TUIs (lazygit-style). |
| Styling | `lipgloss` | Composable styles, theme-friendly. Pairs natively with bubbletea. |
| Components | `bubbles` | Provides `textinput`, `viewport`, `key`, `help` — saves reinventing. |
| Storage (Phase 1) | `modernc.org/sqlite` | Pure Go SQLite, no CGO. Easy cross-compile. |
| Migrations | `golang-migrate/migrate/v4` | File-source migrations, well-documented. |
| Config | `koanf` | Lighter than viper, composable providers. |
| Logging | `log/slog` (stdlib) | No dep, structured, JSON handler built-in. |
| Tests | `testify` + golden files | Familiar assertions; golden snapshots fit TUI rendering. |
| LLM (Phase 2) | Ollama HTTP API | Local-first, swappable models, simple HTTP/SSE. |
| Embeddings (Phase 5) | `nomic-embed-text` via Ollama | Same provider, no extra runtime. |
| Vector store (Phase 5) | SQLite + `sqlite-vec` | Co-locate with task data, hybrid w/ FTS5. |

## Consequences

- **Pure Go everywhere** until Phase 5 vector extension. Cross-compile stays trivial.
- **Ollama dependency** is external (user runs the daemon); we ship only the client. Acceptable since the project is local-first.
- **No ORM**: thin SQL via `database/sql`. Keeps the learning surface clean.
- **Theme system is local-only** in Phase 0 (charm only). Tokyo and mono variants from the design mock are deferred / dropped.

## Alternatives considered

- `tview` instead of `bubbletea`: more widget-heavy but less idiomatic and harder to theme.
- `viper` instead of `koanf`: heavier, more deps. koanf is enough for our YAML+env needs.
- `go-sqlite3` (CGO) instead of `modernc.org/sqlite`: faster but adds CGO toolchain pain.
- LangChainGo / `eino` for agent orchestration: deferred. Phase 6 will build a minimal graph engine for the learning value, then revisit.
