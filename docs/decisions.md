# Locked decisions

These are settled. Don't relitigate without writing an ADR first.

| # | Decision | Why |
|---|----------|-----|
| 1 | **Charm theme only** — peach `#f5a97f`, mauve `#c6a0f6`, warm dark `#16131c`. | Three theme variants from the original mock (charm / tokyo / mono) tripled the surface area for zero learning gain. Refactoring to multi-theme can come later as an exercise. |
| 2 | **English UI strings**, mixed-language seed content allowed. | Keeps grep/eval clean. Code, commits, tests, errors → EN. Personal-project mock data may stay PT for flavor. |
| 3 | **Manual sort order**. `Move(id, delta)` is the only reorder API. | Auto-rank (by priority + due) hides the user's intent. Phase 0 stays explicit; auto-rank may return as an opt-in agent action later. |
| 4 | **Phase 0 = pure TUI**, no persistence, no AI. | Persistence and AI live in their own phases so each one is an isolated learning surface. |
| 5 | **Module path is plain `goz`** (no `github.com/...` prefix yet). | Local-first. Trivial to swap with `go mod edit -module` if/when published. |
| 6 | **Custom graph engine for multi-agent (Phase 6)**, not LangChainGo / Eino / etc. | Building the orchestration layer by hand is a primary learning goal. Comparing to existing libs comes after, not before. |
| 7 | **Pure-Go SQLite (`modernc.org/sqlite`)**, not CGO `go-sqlite3`. | Cross-compile stays trivial; we never need the perf delta at this scale. |
| 8 | **Ollama for both LLM and embeddings**, not OpenAI / Anthropic. | Local-first project. One daemon to run, one HTTP API to learn. |
| 9 | **XDG paths** for data, config, state, cache (Phase 1+). | Standard Linux convention. Backups, dotfile sync, and lazygit-style hygiene all assume this. |
| 10 | **No comments unless the *why* is non-obvious**. Don't restate what the code does. | Code already says what; comments should say why. Echoes the global project guidance in `CLAUDE.md`. |
