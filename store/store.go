package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	dbPath string
	// db connection held here
}

// Open opens or creates the SQLite database in the given data directory
func Open(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("store: creating data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "games.db")

	// TODO open db connection using modernc.org/sqlite
	// enable WAL mode, setup connection pool

	return &Store{
		dbPath: dbPath,
	}, nil
}

// Migrate runs all pending database migration
func (s *Store) Migrate() error {
	// TODO apply embedded SQL migrations in order
	return nil
}

// Close closes db connection
func (s *Store) Close() error {
	// TODO close database connection
	return nil
}

// Game represents a stored chess game
type Game struct {
	ID        int64
	White     string
	Black     string
	WhiteElo  int
	BlackElo  int
	Result    string
	Date      string
	Event     string
	Site      string
	Round     string
	ECO       string
	Opening   string
	PlyCount  int
	Source    string // "played" or "imported" enum?
	Status    string // "in_progress" or "finished" enum?
	Partial   bool
	StartFEN  string
	PGN       string
	Hash      string
	CreatedAt int64
	UpdatedAt int64
}

// GameFilter specifies search criteria for listing games
type GameFilter struct {
	Player   string
	Result   string
	DateFrom string
	DateTo   string
	Source   string
	Analysed *bool // bool pointer seems wrong here
	Limit    int
	Cursor   string
}

// CreateGame inserts a new game and returns the ID
func (s *Store) CreateGame(ctx context.Context, g *Game) (int64, error) {
	// TODO implement
	return 0, nil
}

// UpdateGame updates an existing game.
func (s *Store) UpdateGame(ctx context.Context, g *Game) error {
	// TODO: Implement
	return nil
}

// GetGame retrieves a game by ID.
func (s *Store) GetGame(ctx context.Context, id int64) (*Game, error) {
	// TODO: Implement
	return nil, nil
}

// ListGames returns games matching the filter.
func (s *Store) ListGames(ctx context.Context, f GameFilter) ([]*Game, string, error) {
	// TODO: Implement with FTS5 search and keyset pagination
	// Returns games and next cursor
	return nil, "", nil
}

// DeleteGame removes a game and cascades to analysis.
func (s *Store) DeleteGame(ctx context.Context, id int64) error {
	// TODO: Implement
	return nil
}

// GetInProgressGames returns all games with status "in_progress".
func (s *Store) GetInProgressGames(ctx context.Context) ([]*Game, error) {
	// TODO: Implement
	return nil, nil
}

// AnalysisRun represents a completed analysis of a game.
type AnalysisRun struct {
	ID            int64
	GameID        int64
	EngineName    string
	EngineVersion string
	Settings      string // JSON
	CreatedAt     int64
}

// AnalysisPly represents analysis data for a single ply.
type AnalysisPly struct {
	RunID    int64
	Ply      int
	CP       int
	Mate     int
	BestUCI  string
	BestSAN  string
	Depth    int
	Class    string
	WPLoss   float64
	Accuracy float64
}

// SaveAnalysisRun stores a complete analysis run with all plies.
func (s *Store) SaveAnalysisRun(ctx context.Context, run *AnalysisRun, plies []*AnalysisPly) error {
	// TODO: Implement
	return nil
}

// GetAnalysisRun retrieves the latest analysis run for a game.
func (s *Store) GetAnalysisRun(ctx context.Context, gameID int64) (*AnalysisRun, []*AnalysisPly, error) {
	// TODO: Implement
	return nil, nil, nil
}
