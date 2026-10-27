package sqlitestore

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func (s *Store) Migrate(ctx context.Context) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("sqlitestore: read migrations: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return s.withTx(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version INTEGER PRIMARY KEY,
				name TEXT NOT NULL,
				applied_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))
			)
		`); err != nil {
			return fmt.Errorf("sqlitestore: ensure schema_migrations: %w", err)
		}

		applied, err := appliedMigrations(ctx, tx)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}

			version, err := migrationVersion(entry.Name())
			if err != nil {
				return err
			}

			if applied[version] {
				continue
			}

			path := filepath.Join("migrations", entry.Name())
			body, err := migrationsFS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("sqlitestore: read migration %s: %w", entry.Name(), err)
			}

			if _, err := tx.ExecContext(ctx, string(body)); err != nil {
				return fmt.Errorf("sqlitestore: apply migration %s: %w", entry.Name(), err)
			}

			if _, err := tx.ExecContext(ctx,
				`INSERT INTO schema_migrations(version, name) VALUES (?, ?)`,
				version,
				entry.Name(),
			); err != nil {
				return fmt.Errorf("sqlitestore: record migration %s: %w", entry.Name(), err)
			}
		}

		return nil
	})
}

func appliedMigrations(ctx context.Context, tx *sql.Tx) (map[int]bool, error) {
	rows, err := tx.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("sqlitestore: list applied migrations: %w", err)
	}
	defer rows.Close()

	out := make(map[int]bool)

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("sqlitestore: scan applied migration: %w", err)
		}
		out[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlitestore: iterate applied migrations: %w", err)
	}

	return out, nil
}

func migrationVersion(name string) (int, error) {
	base := filepath.Base(name)
	parts := strings.SplitN(base, "_", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("sqlitestore: invalid migration name %q", name)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("sqlitestore: invalid migration version %q: %w", name, err)
	}

	return version, nil
}
