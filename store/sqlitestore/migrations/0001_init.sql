PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS games (
    id INTEGER PRIMARY KEY,

    white TEXT NOT NULL DEFAULT '',
    black TEXT NOT NULL DEFAULT '',
    white_elo INTEGER,
    black_elo INTEGER,

    result TEXT NOT NULL DEFAULT '*',
    date TEXT NOT NULL DEFAULT '',
    event TEXT NOT NULL DEFAULT '',
    site TEXT NOT NULL DEFAULT '',
    round TEXT NOT NULL DEFAULT '',

    eco TEXT NOT NULL DEFAULT '',
    opening TEXT NOT NULL DEFAULT '',

    ply_count INTEGER NOT NULL DEFAULT 0,

    source TEXT NOT NULL CHECK (source IN ('played', 'imported')),
    status TEXT NOT NULL CHECK (status IN ('in_progress', 'finished')),

    partial INTEGER NOT NULL DEFAULT 0 CHECK (partial IN (0, 1)),

    start_fen TEXT NOT NULL DEFAULT '',
    pgn TEXT NOT NULL,
    hash TEXT NOT NULL UNIQUE,

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE INDEX IF NOT EXISTS idx_games_updated_at_id
    ON games(updated_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_games_date_id
    ON games(date DESC, id DESC);

CREATE INDEX IF NOT EXISTS idx_games_result
    ON games(result);

CREATE INDEX IF NOT EXISTS idx_games_source
    ON games(source);

CREATE INDEX IF NOT EXISTS idx_games_status
    ON games(status);

CREATE INDEX IF NOT EXISTS idx_games_hash
    ON games(hash);

CREATE VIRTUAL TABLE IF NOT EXISTS games_fts
USING fts5(
    white,
    black,
    event,
    site,
    opening,
    eco,
    content='games',
    content_rowid='id'
);

CREATE TRIGGER IF NOT EXISTS games_ai
AFTER INSERT ON games
BEGIN
INSERT INTO games_fts(rowid, white, black, event, site, opening, eco)
VALUES (new.id, new.white, new.black, new.event, new.site, new.opening, new.eco);
END;

CREATE TRIGGER IF NOT EXISTS games_ad
AFTER DELETE ON games
BEGIN
INSERT INTO games_fts(games_fts, rowid, white, black, event, site, opening, eco)
VALUES ('delete', old.id, old.white, old.black, old.event, old.site, old.opening, old.eco);
END;

CREATE TRIGGER IF NOT EXISTS games_au
AFTER UPDATE ON games
BEGIN
INSERT INTO games_fts(games_fts, rowid, white, black, event, site, opening, eco)
VALUES ('delete', old.id, old.white, old.black, old.event, old.site, old.opening, old.eco);

INSERT INTO games_fts(rowid, white, black, event, site, opening, eco)
VALUES (new.id, new.white, new.black, new.event, new.site, new.opening, new.eco);
END;

CREATE TRIGGER IF NOT EXISTS games_set_updated_at
AFTER UPDATE ON games
                            FOR EACH ROW
                            WHEN new.updated_at = old.updated_at
BEGIN
UPDATE games
SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
WHERE id = old.id;
END;

CREATE TABLE IF NOT EXISTS analysis_runs (
    id INTEGER PRIMARY KEY,
    game_id INTEGER NOT NULL,

    engine_id TEXT NOT NULL DEFAULT '',
    engine_name TEXT NOT NULL,
    engine_version TEXT NOT NULL DEFAULT '',
    settings_json TEXT NOT NULL DEFAULT '{}',

    status TEXT NOT NULL CHECK (status IN ('completed', 'cancelled', 'failed')) DEFAULT 'completed',

    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),

    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_analysis_runs_game_created
    ON analysis_runs(game_id, created_at DESC, id DESC);

CREATE TABLE IF NOT EXISTS analysis_plies (
    run_id INTEGER NOT NULL,
    ply INTEGER NOT NULL,

    cp INTEGER,
    mate INTEGER,

    best_uci TEXT NOT NULL DEFAULT '',
    best_san TEXT NOT NULL DEFAULT '',

    depth INTEGER NOT NULL DEFAULT 0,

    class TEXT NOT NULL DEFAULT '',
    wp_loss REAL NOT NULL DEFAULT 0,
    accuracy REAL NOT NULL DEFAULT 100,

    PRIMARY KEY (run_id, ply),

    FOREIGN KEY (run_id) REFERENCES analysis_runs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS engines_cache (
    engine_id TEXT PRIMARY KEY,
    engine_name TEXT NOT NULL DEFAULT '',
    engine_version TEXT NOT NULL DEFAULT '',
    options_json TEXT NOT NULL DEFAULT '{}',
    probed_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
);
