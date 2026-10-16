# Spec: Local Chess Interface (Go)

**Date:** 2026-10-15
**Status:** Draft
**Author:** Leo

---

## 1. Summary

A single-binary Go application, run locally, that serves a browser UI on `127.0.0.1`. It lets the user (a) play games against locally installed UCI engines such as Stockfish or Lc0, (b) import games from PGN files and run a full engine review of them (per-move evaluation, best move, mistake classification, accuracy, eval graph, plus live analysis of any position), and (c) keep every played or imported game in a local, searchable library backed by SQLite.

The UI is server-rendered with htmx. All chess logic lives in Go. There is no Node toolchain and no cgo.

---

## 2. Background & Motivation

- Existing GUIs (Arena, Lichess, Chessbase, Scid) are either heavy, cloud-bound, or awkward to script and package. This project wants a small, self-contained, hackable tool that fits a homelab / NixOS workflow: one binary, one data directory, no external services.
- Engines are already installed locally and speak UCI. The app should treat them as plain subprocesses configured by path.
- Prior art to review before implementation: Scid vs. PC, Lucas Chess, CuteChess (engine matches, UCI handling), Lichess analysis (move classification and accuracy approach).

---

## 3. Goals

- [ ] Play a game against any configured UCI engine, as white, black or random, with fixed move time or depth, and an optional Elo limit where the engine supports it.
- [ ] Persist every game, including in-progress ones, so a crash or restart never loses a game and a game can be resumed.
- [ ] Import one or many games from PGN (files or pasted text), tolerating real-world messy PGN, de-duplicating on re-import.
- [ ] Run a full-game engine review: per-ply eval, best move, move classification (best / good / inaccuracy / mistake / blunder), per-side accuracy, eval graph. Results are stored and reloaded without re-running the engine.
- [ ] Live analysis of the currently displayed position with multi-PV.
- [ ] Searchable library by player, result, date, event, opening/ECO, source, analysed status.
- [ ] Ship as one binary with embedded templates and static assets; data in an XDG-style directory.
- [ ] Rules correctness verified by perft tests; UCI client verified against a fake engine.

---

## 4. Non-Goals

- Engine protocols other than UCI (no XBoard/WinBoard).
- Engine-vs-engine matches, tournaments, or Elo rating of engines.
- Position or opening explorer across the library (searching by FEN / move sequence).
- Opening-book or tablebase integration; ECO naming beyond using PGN headers if present.
- Multi-user access, accounts, remote access, or any network exposure beyond loopback.
- Online play, Lichess/Chess.com import, cloud sync.
- Chess variants (Chess960 is deferred; see Open Questions).
- Real chess clocks with increment (`wtime/btime`) in v1; v1 uses fixed movetime / depth / nodes.
- A JS framework or build step.
- Scale beyond roughly 100k games in the library.

---

## 5. Proposed Solution

### 5.1 Overview

One Go process:

1. **HTTP server** (`net/http`, Go 1.22+ pattern routing, `html/template`, `go:embed`) bound to `127.0.0.1`. Pages are full HTML; interactions are htmx partial swaps.
2. **Engine manager** owns UCI subprocesses. Play sessions and analysis jobs each use their own process so they never contend.
3. **Store** (SQLite via `modernc.org/sqlite`, pure Go) holds games, analysis, and settings.
4. **Chess core**: a maintained third-party Go chess library for move generation, FEN, SAN and PGN, wrapped in a thin internal adapter so it can be replaced.
5. **Background workers** run analysis jobs and report progress via Server-Sent Events (htmx SSE extension).

### 5.2 Design Details

#### 5.2.1 Package layout

```
cmd/localchess/        main, flags, wiring, graceful shutdown
internal/chessx/       adapter over the chess library (Position, Move, SAN, FEN, outcome)
internal/pgnio/        import/export, normalisation, dedupe hashing, tolerant parsing
internal/uci/          UCI client (process, protocol, info/bestmove parsing)
internal/engines/      engine registry, config file, PATH auto-detect, option handling
internal/play/         play sessions, engine-turn orchestration, autosave
internal/analysis/     job runner, evaluation pipeline, classification, accuracy
internal/store/        SQLite, migrations, queries, FTS
internal/web/          handlers, templates, SSE, embedded static assets, security middleware
internal/web/static/   htmx.min.js (vendored), piece SVG sprite, CSS
testdata/              PGN corpus, perft suites
```

#### 5.2.2 Chess core

- Use an existing library for legal move generation, FEN, SAN and PGN. ⚠️ Candidates (e.g. `notnil/chess` and its maintained forks) must be evaluated before commit on: maintenance status, PGN tolerance (comments, variations, NAGs, malformed moves), SAN round-trip correctness, perft speed, and licence. Decision is a Milestone 1 exit criterion.
- Evaluation criteria also include: preserving variations, comments, NAGs and `[%cal]` / `[%csl]` annotations on parse and re-export without loss (needed for the future iterations in section 11).
- All library types stay behind `internal/chessx`. No other package imports the library directly. This keeps a swap or a partial home-grown replacement (e.g. a custom PGN reader) cheap.
- `chessx` exposes a move-tree type (node = move, comment, NAGs, annotations, children). A game's mainline is one view of that tree. v1 only uses the mainline, but the type is tree-shaped from the start so study mode needs no rewrite.
- The stored PGN text is the source of truth for a game. Anything derived (mainline plies, analysis rows) can be rebuilt from it, so later features need no data migration.

#### 5.2.3 UCI client (`internal/uci`)

A thin, own implementation (not the library's bundled UCI helper, ⚠️ confirm after evaluation) so process lifecycle is fully controlled.

```go
type Engine struct { /* cmd, stdin, stdout, reader goroutine, state */ }

func Start(ctx context.Context, path string, args []string) (*Engine, error)
func (e *Engine) Options() map[string]Option           // from `uci` handshake
func (e *Engine) SetOption(name, value string) error
func (e *Engine) NewGame(ctx context.Context) error    // ucinewgame + isready
func (e *Engine) Go(ctx context.Context, pos Position, lim Limits, opts GoOpts) (<-chan Info, error)
func (e *Engine) Stop(ctx context.Context) error
func (e *Engine) Close() error                         // quit, then kill after grace period
```

- `Limits`: movetime, depth, nodes. `GoOpts`: MultiPV, searchmoves (unused v1).
- `Info` carries depth, seldepth, multipv, score (cp or mate, with bound flags), nodes, nps, time, pv. A final `BestMove` event ends the stream.
- Handshake: `uci` → collect `id` / `option` lines → `uciok`; `isready` → `readyok` with timeout.
- Context cancellation sends `stop`, waits for `bestmove` with a grace timeout, then kills the process. Reader goroutines exit on EOF; `Close` is idempotent. No goroutine may outlive `Close`.
- Crash handling: unexpected process exit surfaces as a typed error (`ErrEngineExited`), the session is marked engine-failed, and the UI offers restart. Engine stderr is captured to a bounded ring buffer for diagnostics.
- Errors wrapped with `fmt.Errorf("uci: %w", err)`; sentinel errors for timeout, exit, and unsupported option.
- Unknown options are ignored with a logged warning, never fatal.

#### 5.2.4 Engine registry (`internal/engines`)

- Config file `engines.toml` (or JSON, ⚠️ pick one) in the config dir: `id`, `name`, `path`, `args`, default `options`.
- Auto-detect common binaries on `PATH` (stockfish, lc0, and similar); detected entries are shown but not persisted until the user saves them.
- Adding an engine launches it once to run the handshake and cache advertised options (UCI_Elo range, Skill Level, Threads, Hash, MultiPV).
- Elo limiting: if `UCI_LimitStrength` and `UCI_Elo` exist, clamp the requested Elo into the advertised range and set both. Otherwise the UI disables the field and says why.
- ⚠️ Security: engine paths are executed. Configuration endpoints require the origin protections in 5.2.8.

#### 5.2.5 Play sessions (`internal/play`)

- A session = one game, one engine process, one `chessx.Position`, move list, status.
- Interaction is fully server-driven (htmx): each square is an element that posts `select`/`move` requests; the server keeps selection state and returns a board fragment with legal-target highlights. Promotion shows a server-rendered chooser fragment. Localhost latency makes this acceptable, ⚠️ validate feel in Milestone 3.
- Engine turn: after a user move the response returns immediately with the updated board and an "engine thinking" fragment that opens an SSE stream; the engine's reply and final board arrive over SSE. Undo is rejected while the engine is thinking, or cancels the search first (decision: reject in v1).
- Game end: checkmate, stalemate, insufficient material, fivefold, 75-move automatic; threefold and 50-move claimable. v1 auto-claims threefold and 50-move for simplicity. Resign supported. Draw offers not supported.
- **Autosave**: the game row is created at game start with `result='*'` and `status='in_progress'`, and updated after every ply. Finished games flip to a final result. On startup, in-progress games are listed as "Resume".
- Engine process closed on game end, new game, or shutdown.

#### 5.2.6 Import and library

- **Import** accepts multipart file upload (multiple files) and pasted text. Parsing streams game by game; the request returns a result summary (added / duplicates / skipped / partial with reasons).
- Tolerance rules: ignore unknown tags, accept missing seven-tag-roster fields, keep comments and NAGs in the stored PGN, import the mainline for analysis, truncate at the first illegal move and flag the game `partial`. Games with zero moves are skipped.
- **De-duplication**: hash of normalised key headers (White, Black, Date, Event, Round) + mainline move sequence + starting FEN. Played games include a unique component so repeated engine games never collide.
- **Search**: SQLite FTS5 over players, event, site, opening (⚠️ confirm FTS5 is available in the chosen pure-Go driver build), plus indexed filters on result, date range, source, and analysed flag. Pagination via keyset.
- Export: single game or filtered set as PGN.
- Deletion cascades analysis.

#### 5.2.7 Analysis (`internal/analysis`)

Pipeline per game:

1. Replay mainline into positions `P0..Pn`.
2. For each position run `go` with the configured movetime (default 500 ms) or depth, MultiPV=1. Terminal positions are scored directly (mate / draw) without calling the engine.
3. Store per position: white-POV cp (mate mapped to a large sentinel), mate distance, best move UCI/SAN, depth.
4. Per move `i`, compute mover-POV win probability before and after, using a logistic cp→win% curve with cp clipped to ±1000. The win-probability drop gives the class, and accuracy follows from the drop.
5. Class thresholds (starting values, tunable constants): drop ≥ 15 blunder, ≥ 10 mistake, ≥ 5 inaccuracy; played move equal to best → best; otherwise good. ⚠️ Constants are inspired by Lichess' published approach; verify and tune against sample games.
6. Persist run metadata (engine, version string, settings) and per-ply rows. Re-running replaces the prior run for the same game.

Execution:

- Jobs run in a bounded worker pool (default 1 concurrent job). Each job owns its engine process and a `context.Context`; cancel stops the engine and marks the run cancelled with partial results discarded.
- Progress is published over SSE (`done/total`). Closing the browser tab does not cancel the job.
- **Live analysis**: a separate long-lived analysis engine per chosen engine id, guarded by a mutex; new request cancels the previous search (latest-wins). Output streams as SSE lines (score, depth, SAN PV) for up to MultiPV lines.

Board fragments include an empty SVG overlay layer above the squares, reserved for arrows and square highlights in future iterations (section 11).

Review UI: board, move list with classification marks, eval bar, eval graph (server-rendered inline SVG, clickable via htmx to jump to a ply), per-side accuracy summary, "better move" hint on mistakes, keyboard navigation (small amount of vanilla JS only for key bindings).

#### 5.2.8 Web layer and security

- Bind to `127.0.0.1` only by default; `--addr` flag to change, with a warning when non-loopback.
- Middleware rejects requests whose `Host` is not loopback (DNS-rebinding defence) and any state-changing request whose `Origin`/`Referer` does not match the host. Per-process CSRF token embedded in forms/htmx headers.
- Rationale: engine configuration executes arbitrary binaries, so a hostile web page must not be able to drive the local API.
- Static assets (htmx, SVG sprite, CSS) vendored and `go:embed`ded. No CDN calls; the app works offline.
- Graceful shutdown: stop accepting requests, cancel jobs, close engines, flush SQLite.

#### 5.2.9 Data model (initial)

```sql
games(id, white, black, white_elo, black_elo, result, date, event, site, round,
      eco, opening, ply_count, source,            -- 'played' | 'imported'
      status,                                      -- 'in_progress' | 'finished'
      partial, start_fen, pgn, hash UNIQUE, created_at, updated_at)

games_fts(white, black, event, site, opening)     -- FTS5, content=games

analysis_runs(id, game_id FK, engine_name, engine_version, settings_json, created_at)
analysis_plies(run_id FK, ply, cp, mate, best_uci, best_san, depth,
               class, wp_loss, accuracy)          -- ply 0 = start position

engines_cache(engine_id, options_json, probed_at)
settings(key, value)
schema_migrations / PRAGMA user_version
```

- Migrations are embedded SQL files applied in order inside a transaction at startup.
- WAL mode; one write connection, a small read pool.

#### 5.2.10 Configuration and files

| Item | Default location | Override |
|------|------------------|----------|
| Database | `$XDG_DATA_HOME/localchess/games.db` | `--data-dir` / `LOCALCHESS_DATA` |
| Engine config | `$XDG_CONFIG_HOME/localchess/engines.toml` | `--config-dir` |
| Listen address | `127.0.0.1:8080` | `--addr` |

### 5.3 Affected Components

New project, no existing code is modified.

| Component | Change Type | Notes |
|-----------|-------------|-------|
| `cmd/localchess` | New | Entry point, wiring, shutdown |
| `internal/chessx` | New | Adapter over chosen chess library |
| `internal/pgnio` | New | Tolerant import/export, hashing |
| `internal/uci` | New | UCI client |
| `internal/engines` | New | Registry, config, detection |
| `internal/play` | New | Sessions and autosave |
| `internal/analysis` | New | Job runner, classification |
| `internal/store` | New | SQLite and migrations |
| `internal/web` | New | Handlers, templates, SSE, middleware |
| Packaging | New | `go build`; Nix flake recommended for this environment (see Open Questions) |

---

## 6. Alternatives Considered

| Option | Pros | Cons | Reason Not Chosen |
|--------|------|------|-------------------|
| TUI (Bubble Tea) | No browser, very light | Poor fit for board, eval graph, game list | UI choice: embedded web |
| Native desktop (Fyne/Gio) | Real app feel | Heavy toolchain, packaging friction | UI choice: embedded web |
| SPA (TypeScript/Svelte) | Smoother board interaction, local move validation | Adds Node build to a pure-Go project | Frontend choice: htmx |
| Vanilla JS + JSON API | Instant click-to-move, simple server | Chess state duplicated client-side or fetch per move | Frontend choice: htmx (user decision) |
| cgo SQLite (`mattn/go-sqlite3`) | Mature, fast | Needs C toolchain, harder on NixOS and cross-builds | Pure Go driver preferred |
| Flat PGN files + custom index | Human-readable | Search and analysis caching need bespoke indexing | SQLite chosen |
| Own move generator | No dependency | Large correctness surface | Existing library chosen |
| Library's bundled UCI helper | Less code | Less control over lifecycle, cancellation, crash handling | Thin own client preferred ⚠️ revisit after evaluation |

---

## 7. Acceptance Criteria

- [ ] Perft suite passes for the standard start position and the well-known tricky positions (castling rights, en passant, promotion, pins) to a documented depth.
- [ ] UCI client passes its fake-engine tests: handshake, options, go/info/bestmove, stop, timeout, crash, malformed output.
- [ ] A full game can be played against a real engine (manual check with Stockfish) as white and black, including promotion, undo, resign and each game-end condition.
- [ ] Killing the process mid-game and restarting offers to resume that game at the correct position.
- [ ] Importing a multi-game PGN with comments, variations, NAGs and one malformed game adds the valid games, flags the malformed one `partial`, and re-importing the same file adds zero.
- [ ] Full-game analysis of a sample game stores results, survives restart, and renders graph, classifications and accuracy without re-running the engine.
- [ ] Analysis can be cancelled; no engine processes or goroutines remain afterwards (checked under `-race` with a leak check).
- [ ] Library search filters by player, result, date, source and analysed flag, and returns in under 100 ms on a 50k-game library.
- [ ] Requests with a non-loopback `Host` or cross-origin state-changing requests are rejected.
- [ ] `go build` produces a single binary that works with no network access.

---

## 8. Test Plan

| Test Type | Scenario | Expected Result |
|-----------|----------|-----------------|
| Unit (table-driven) | Perft counts for start position, Kiwipete, and other standard suites | Counts match published values |
| Unit | SAN ↔ UCI round-trip on tricky moves (disambiguation, castling, promotion) | Exact match |
| Unit | UCI `info` line parsing (cp, mate, bounds, multipv, pv, truncated lines) | Correct `Info`, no panic on garbage |
| Integration, fake engine | Handshake, option set, `go movetime`, `bestmove` | Moves returned, options recorded |
| Integration, fake engine | Engine hangs after `go` | `stop` sent, killed after grace, typed timeout error |
| Integration, fake engine | Engine exits mid-search | `ErrEngineExited`, session marked failed, no leaked goroutines |
| Integration, fake engine | Engine ignores `isready` | Handshake timeout error |
| Unit | Win% / accuracy / classification on fixed cp sequences | Matches expected table |
| Integration | PGN corpus import (comments, variations, NAGs, bad headers, illegal move, multi-GB-style repeated games) | Counts and flags as specified |
| Integration | Re-import same file | 0 added, N duplicates |
| Integration | Analysis job on a short game with fake engine returning scripted scores | Stored plies and classes match script |
| Integration | Cancel analysis mid-run | Job cancelled, process gone |
| Integration (`httptest`) | Play flow: new game, select, move, engine reply via SSE, resign | Correct fragments and final state |
| Integration (`httptest`) | Host/Origin/CSRF rejection | 403 |
| Store | Migrations from empty DB; FTS search; cascade delete | Expected rows |
| Race / leak | Full suite with `-race` and goroutine-leak detection | Clean |
| Manual | Play vs Stockfish and Lc0, Elo-limited play, review a real tournament game | Behaves as expected |

Fake engine approach: the test binary re-executes itself as a scripted UCI engine (selected by an environment variable), so tests need no external engine and no extra build step. ⚠️ Optional later: opt-in real-Stockfish tests behind a build tag (not an acceptance gate).

---

## 9. Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Chosen chess library has PGN or rules gaps | Medium | High | Evaluate in Milestone 1 with perft and PGN corpus; keep it behind `chessx` |
| htmx click-to-move feels laggy or fiddly (promotion, drag) | Medium | Medium | Server-side selection state; validate early in Milestone 3; allow a small JS enhancement for drag only if needed |
| UCI engine quirks (slow `uciok`, non-compliant info lines) | High | Medium | Tolerant parser, timeouts, fake-engine tests, test against Stockfish and Lc0 |
| Zombie engine processes / goroutine leaks | Medium | High | Context-driven lifecycle, `Close` idempotent, kill after grace, leak checks in CI |
| Hostile web page drives the local API (engine path is executable) | Low | High | Loopback bind, Host/Origin checks, CSRF token, no CORS |
| Pure-Go SQLite lacks FTS5 or is slow on large imports | Low | Medium | Verify early; batch imports in one transaction; fall back to indexed LIKE |
| Classification thresholds feel wrong | Medium | Low | Constants isolated and unit-tested; tune on sample games |
| Large PGN uploads exhaust memory | Medium | Medium | Stream parse game by game; cap upload size |

---

## 10. Rollout Plan

- Feature flag: not applicable (new application).
- Migration: not applicable; schema migrations are embedded from the first version.
- Rollback: not applicable; data directory is a single SQLite file plus config, easy to back up.
- Suggested milestones (⚠️ rough effort, single developer):
  1. **Foundations**: chess library evaluation (including lossless PGN annotation round-trip), tree-shaped `chessx`, perft tests, `uci` client plus fake engine tests (~1 week).
  2. **Store and import**: schema, migrations, PGN import/export, dedupe, library list and search (~1 week).
  3. **Play**: web shell, board fragments, play sessions, autosave and resume (~1 week).
  4. **Analysis**: job runner, classification, eval graph, live analysis, SSE (~1 to 1.5 weeks).
  5. **Hardening**: security middleware, error states, packaging, docs (~3 to 4 days).

---

## 11. Future Iterations (out of scope for v1)

These are deliberately not built in v1. The v1 design includes the small hooks listed under each so they can be added without rework.

### 11.1 Board arrows

- **Engine top-move arrows**: during live analysis, draw an arrow for each MultiPV line's first move, rendered server-side into the SVG overlay layer. Width or colour encodes rank (best move strongest). Toggle on/off per session.
- **User-drawn arrows and highlights**: right-drag draws an arrow, right-click highlights a square. This needs a small amount of vanilla JS (htmx cannot handle drags); the JS posts the resulting arrow set to the server, which re-renders the overlay.
- **Persistence**: arrows and highlights are stored in PGN comments using the `[%cal ...]` (arrows) and `[%csl ...]` (squares) convention used by Lichess and ChessBase, so they survive export and import.
- **v1 hooks**: SVG overlay layer in board fragments; annotation preservation as a chess library evaluation criterion; `chessx` node type carries annotations.

### 11.2 Study mode

- **Purpose**: study chess books and courses by drawing out several lines in one session and comparing candidate moves.
- **Model**: a *study* contains *chapters* (e.g. one per book chapter or course lesson). Each chapter holds a move tree: nodes with a move, comment, NAGs and arrows, and arbitrary branching. Chapters import from and export to PGN with variations.
- **Candidate comparison**: from any node, show all child moves side by side, each with its engine eval (cached per position), with one arrow per candidate on the board.
- **Navigation and editing**: add a move creates a child (or reuses an existing one), promote a variation to main line, delete a branch, edit comments and NAGs.
- **Transpositions**: ⚠️ the same position can be reached by different move orders. Decide whether nodes are a pure tree (duplicates allowed, simplest) or a position graph (merged, more complex). Recommend a pure tree first, with a "this position also appears at ..." hint derived from FEN lookup.
- **Storage**: `studies(id, title, notes, created_at, updated_at)` and `study_chapters(id, study_id, title, position, pgn)`. Chapter PGN is the source of truth; parsed trees are rebuilt on load. Engine evals reuse the existing per-position analysis cache keyed by FEN and engine settings.
- **Search**: study titles and chapter titles join the FTS index.
- **v1 hooks**: tree-shaped `chessx` type; stored-PGN-as-source-of-truth; analysis cache keyed by position rather than only by game ply.

### 11.3 Risks specific to these features

| Risk | Mitigation |
|------|------------|
| Tree editing with htmx becomes clunky | Server holds the cursor (current node); every edit is a small fragment swap; prototype early |
| Annotation round-trip loses data in the chosen library | Make lossless round-trip a library selection gate in Milestone 1 |
| Eval cache grows large across studies | Cache by FEN with a size cap and eviction, or store only at user-requested depth |

---

## 12. Open Questions

1. Which chess library? Evaluate candidates in Milestone 1 (maintenance, PGN tolerance, perft speed, licence).
2. Engine config format: TOML or JSON? (Another dependency vs. stdlib only.)
3. Packaging: plain `go build` only, or also a Nix flake providing the package and a dev shell?
4. Does the UI need real clocks (`wtime/btime`, increment) soon, or is fixed movetime/depth enough for v1?
5. Chess960 support: defer, or design the position and PGN layers to allow it from the start?
6. Should draw offers and engine-initiated resignation be handled, or stay out?
7. Undo during engine thinking: reject (current plan) or cancel the search and undo?
8. Should analysis default to depth-based or movetime-based limits, and what defaults (500 ms per position)?
9. Do variations and comments need to be displayed in Review, or only preserved in the stored PGN?
10. Real-Stockfish integration tests: add as an optional tagged suite?
11. Study mode: pure move tree or merged position graph for transpositions?
12. Should arrow and highlight annotations be editable in Review mode as well as study mode?

---

## 13. References

- UCI protocol specification (engine interface documentation)
- Lichess documentation on accuracy and move classification
- Chess Programming Wiki: perft results and test positions
- htmx documentation, including the SSE extension
- `modernc.org/sqlite` documentation (FTS5 support, WAL)
