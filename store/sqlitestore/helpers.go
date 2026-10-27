package sqlitestore

import (
	"database/sql"
	"encoding/base64"
	"encoding/json/v2"
	"errors"
	"fmt"
	"neochess/store"
	"time"

	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func isConstraintError(err error) bool {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		code := sqliteErr.Code()
		return code == sqlitelib.SQLITE_CONSTRAINT ||
			code == sqlitelib.SQLITE_CONSTRAINT_UNIQUE ||
			code == sqlitelib.SQLITE_CONSTRAINT_PRIMARYKEY
	}

	return false
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func defaultString(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

func parseDBTime(v string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.000Z", v)
}

func formatDBTime(v time.Time) string {
	return v.UTC().Format("2006-01-02T15:04:05.000Z")
}

type gameCursor struct {
	UpdatedAt string `json:"updated_at"`
	ID        int64  `json:"id"`
}

func encodeGameCursor(updatedAt time.Time, id int64) string {
	body, _ := json.Marshal(gameCursor{
		UpdatedAt: formatDBTime(updatedAt),
		ID:        id,
	})

	return base64.RawURLEncoding.EncodeToString(body)
}

func decodeGameCursor(cursor string) (time.Time, int64, error) {
	if cursor == "" {
		return time.Time{}, 0, nil
	}

	body, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("sqlitestore: decode cursor: %w", err)
	}

	var c gameCursor
	if err := json.Unmarshal(body, &c); err != nil {
		return time.Time{}, 0, fmt.Errorf("sqlitestore: unmarshal cursor: %w", err)
	}

	t, err := parseDBTime(c.UpdatedAt)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("sqlitestore: parse cursor time: %w", err)
	}

	return t, c.ID, nil
}

func nullIntPtr(v *int) any {
	if v == nil {
		return nil
	}

	return *v
}

func intPtrFromNull(v sql.NullInt64) *int {
	if !v.Valid {
		return nil
	}

	out := int(v.Int64)
	return &out
}

func sqlErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return store.ErrNotFound
	}

	if isConstraintError(err) {
		return store.ErrDuplicate
	}

	return err
}
