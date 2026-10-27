package sqlitestore

import (
	"context"
	"database/sql"
	"fmt"
	"neochess/store"
	"os"
	"path/filepath"
)

var _ store.Store = (*Store)(nil)

type Store struct {
	db *sql.DB
}

type Options struct {
	DataDir string
	DBPath  string
}

func Open(ctx context.Context, opts Options) (*Store, error) {
	dbPath := opts.DBPath
	if dbPath == "" {
		if opts.DataDir == "" {
			return nil, fmt.Errorf("sqlitestore: data dir or db path is required")
		}

		if err := os.MkdirAll(opts.DataDir, 0o755); err != nil {
			return nil, fmt.Errorf("sqlitestore: create data dir: %w", err)
		}

		dbPath = filepath.Join(opts.DataDir, "games.db")
	}

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlitestore: open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	s := &Store{db: db}

	if err := s.configure(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := s.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return s, nil
}

func (s *Store) configure(ctx context.Context) error {
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA busy_timeout = 5000",
		"PRAGMA temp_store = MEMORY",
	}

	for _, pragma := range pragmas {
		if _, err := s.db.ExecContext(ctx, pragma); err != nil {
			return fmt.Errorf("sqlitestore: configure sqlite %q: %w", pragma, err)
		}
	}

	return nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("sqlitestore: ping: %w", err)
	}

	return nil
}

func (s *Store) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sqlitestore: begin tx: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("sqlitestore: commit tx: %w", err)
	}

	return nil
}
