package sqlitestore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"neochess/store"
	"strings"
)

func (s *Store) CreateGame(ctx context.Context, params store.CreateGameParams) (int64, error) {
	result := defaultString(params.Result, "*")
	status := params.Status
	if status == "" {
		status = store.GameStatusFinished
	}

	var id int64

	err := s.db.QueryRowContext(ctx, `
		INSERT INTO games (
			white, black, white_elo, black_elo,
			result, date, event, site, round,
			eco, opening, ply_count,
			source, status, partial,
			start_fen, pgn, hash
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
		`,
		params.White,
		params.Black,
		nullIntPtr(params.WhiteElo),
		nullIntPtr(params.BlackElo),
		result,
		params.Date,
		params.Event,
		params.Site,
		params.Round,
		params.ECO,
		params.Opening,
		params.PlyCount,
		params.Source,
		status,
		boolInt(params.Partial),
		params.StartFEN,
		params.PGN,
		params.Hash,
	).Scan(&id)

	if err != nil {
		if isConstraintError(err) {
			return 0, store.ErrDuplicate
		}
		return 0, fmt.Errorf("sqlite: create game %w", err)
	}
	return id, nil
}

func (s *Store) InsertImportedGame(ctx context.Context, params store.CreateGameParams) (store.ImportGameResult, error) {
	params.Source = store.GameSourceImported

	id, err := s.CreateGame(ctx, params)
	if err == nil {
		return store.ImportGameResult{
			GameID:  id,
			Created: true,
		}, nil
	}

	if !errors.Is(err, store.ErrDuplicate) {
		return store.ImportGameResult{}, err
	}

	existingID, err := s.GetGameIDByHash(ctx, params.Hash)
	if err != nil {
		return store.ImportGameResult{}, err
	}

	return store.ImportGameResult{
		GameID:  existingID,
		Created: false,
	}, nil
}

func (s *Store) GetGame(ctx context.Context, id int64) (store.Game, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT
		    id,
		    white, black, white_elo, black_elo,
		    result, date, event, site, round,
		    eco, opening, ply_count,
		    source, status, partial,
		    start_fen, pgn, hash,
		    created_at, updated_at
		FROM games
		WHERE id = ?
	`, id)

	g, err := scanGame(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.Game{}, store.ErrNotFound
		}
		return store.Game{}, fmt.Errorf("sqlite: get game: %w", err)
	}

	return g, nil
}

func (s *Store) GetGameIDByHash(ctx context.Context, hash string) (int64, error) {
	var id int64

	err := s.db.QueryRowContext(ctx, `
		SELECT id
		FROM games
		WHERE hash = ?
	`, hash).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, store.ErrNotFound
		}
		return 0, fmt.Errorf("sqlite: get game id by hash: %w", err)
	}
	return id, nil
}

func (s *Store) ListGames(ctx context.Context, filter store.GameListFilter) (store.GameListPage, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	args := make([]any, 0)
	where := make([]string, 0)

	from := "games g"

	if strings.TrimSpace(filter.Query) != "" {
		from = "games g JOIN games_fts fts ON fts.rowid = g.id"
		where = append(where, "games_fts MATCH ?")
		args = append(args, filter.Query)
	}

	if filter.Player != "" {
		where = append(where, "(g.white = ? OR g.black = ?)")
		args = append(args, filter.Player, filter.Player)
	}

	if filter.Result != "" {
		where = append(where, "g.result = ?")
		args = append(args, filter.Result)
	}

	if filter.Source != "" {
		where = append(where, "g.source = ?")
		args = append(args, filter.Source)
	}

	if filter.Status != "" {
		where = append(where, "g.status = ?")
		args = append(args, filter.Status)
	}

	if filter.Partial != nil {
		where = append(where, "g.partial = ?")
		args = append(args, boolInt(*filter.Partial))
	}

	if filter.DateFrom != "" {
		where = append(where, "g.date >= ?")
		args = append(args, filter.DateFrom)
	}

	if filter.DateTo != "" {
		where = append(where, "g.date <= ?")
		args = append(args, filter.DateTo)
	}

	if filter.Analysed != nil {
		if *filter.Analysed {
			where = append(where, "EXISTS (SELECT 1 FROM analysis_runs ar WHERE ar.game_id = g.id AND ar.status = 'completed')")
		} else {
			where = append(where, "NOT EXISTS (SELECT 1 FROM analysis_runs ar WHERE ar.game_id = g.id AND ar.status = 'completed')")
		}
	}

	if filter.Cursor != "" {
		afterUpdateAt, afterID, err := decodeGameCursor(filter.Cursor)
		if err != nil {
			return store.GameListPage{}, err
		}

		where = append(where, "(g.updated_at < ? OR (g.updated_at = ? AND g.id < ?))")
		cursorTime := formatDBTime(afterUpdateAt)
		args = append(args, cursorTime, cursorTime, afterID)
	}

	query := `
		SELECT
			g.id,
			g.white, g.black, g.white_elo, g.black_elo,
			g.result, g.date, g.event, g.site, g.round,
			g.eco, g.opening, g.ply_count,
			g.source, g.status, g.partial,
			g.start_fen, g.pgn, g.hash,
			g.created_at, g.updated_at
		FROM ` + from

	if len(where) > 0 {
		query += "\nWHERE " + strings.Join(where, " AND ")
	}

	query += "\nORDER BY g.updated_at DESC, g.id DESC\nLIMIT ?"
	args = append(args, limit+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return store.GameListPage{}, fmt.Errorf("sqlite: list games: %w", err)
	}
	defer rows.Close()

	games, err := scanGames(rows)
	if err != nil {
		return store.GameListPage{}, err
	}

	var nextCursor string
	if len(games) > limit {
		last := games[limit-1]
		nextCursor = encodeGameCursor(last.UpdatedAt, last.ID)
		games = games[:limit]
	}

	return store.GameListPage{
		Games:      games,
		NextCursor: nextCursor,
	}, nil
}

func (s *Store) ListInProgressGames(ctx context.Context) ([]store.Game, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id,
			white, black, white_elo, black_elo,
			result, date, event, site, round,
			eco, opening, ply_count,
			source, status, partial,
			start_fen, pgn, hash,
			created_at, updated_at
		FROM games
		WHERE status = 'in_progress'
		ORDER BY updated_at DESC, id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("sqlitestore: list in-progress games: %w", err)
	}
	defer rows.Close()

	return scanGames(rows)
}

func (s *Store) UpdateGame(ctx context.Context, params store.UpdateGameParams) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE games
		SET
		    white = ?,
			black = ?,
			white_elo = ?,
			black_elo = ?,
			result = ?,
			date = ?,
			event = ?,
			site = ?,
			round = ?,
			eco = ?,
			opening = ?,
			ply_count = ?,
			source = ?,
			status = ?,
			partial = ?,
			start_fen = ?,
			pgn = ?,
			hash = ?
		WHERE id = ?
	`,
		params.White,
		params.Black,
		nullIntPtr(params.WhiteElo),
		nullIntPtr(params.BlackElo),
		defaultString(params.Result, "*"),
		params.Date,
		params.Event,
		params.Site,
		params.Round,
		params.ECO,
		params.Opening,
		params.PlyCount,
		params.Source,
		params.Status,
		boolInt(params.Partial),
		params.StartFEN,
		params.PGN,
		params.Hash,
		params.ID,
	)

	if err != nil {
		if isConstraintError(err) {
			return store.ErrDuplicate
		}
		return fmt.Errorf("sqlite: update game: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite: update game rows affected: %w", err)
	}

	if n == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *Store) DeleteGame(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM games WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("sqlite: delete game: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite: delete game rows affected: %w", err)
	}

	if n == 0 {
		return store.ErrNotFound
	}

	return nil
}

func scanGames(rows *sql.Rows) ([]store.Game, error) {
	var games []store.Game

	for rows.Next() {
		g, err := scanGame(rows)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite: iterate games: %w", err)
	}

	return games, nil
}

func scanGame(row rowScanner) (store.Game, error) {
	var g store.Game
	var whiteElo sql.NullInt64
	var blackElo sql.NullInt64
	var createdAt string
	var updatedAt string
	var partial int

	err := row.Scan(
		&g.ID,
		&g.White,
		&g.Black,
		&whiteElo,
		&blackElo,
		&g.Result,
		&g.Date,
		&g.Event,
		&g.Site,
		&g.Round,
		&g.ECO,
		&g.Opening,
		&g.PlyCount,
		&g.Source,
		&g.Status,
		&partial,
		&g.StartFEN,
		&g.PGN,
		&g.Hash,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return store.Game{}, err
	}

	g.WhiteElo = intPtrFromNull(whiteElo)
	g.BlackElo = intPtrFromNull(blackElo)
	g.Partial = partial != 0

	g.CreatedAt, err = parseDBTime(createdAt)
	if err != nil {
		return store.Game{}, fmt.Errorf("sqlite: parse game created_at: %w", err)
	}

	g.UpdatedAt, err = parseDBTime(updatedAt)
	if err != nil {
		return store.Game{}, fmt.Errorf("sqlite: parse game updated_at: %w", err)
	}

	return g, nil
}
