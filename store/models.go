package store

import "time"

type GameSource string

const (
	GameSourcePlayed   GameSource = "played"
	GameSourceImported GameSource = "imported"
)

type GameStatus string

const (
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusFinished   GameStatus = "finished"
)

// Game represents a stored chess game
type Game struct {
	ID int64

	White    string
	Black    string
	WhiteElo *int
	BlackElo *int

	Result string
	Date   string
	Event  string
	Site   string
	Round  string

	ECO     string
	Opening string

	PlyCount int

	Source GameSource // "played" or "imported" enum?
	Status GameStatus // "in_progress" or "finished" enum?

	Partial bool

	StartFEN string
	PGN      string
	Hash     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateGameParams struct {
	White    string
	Black    string
	WhiteElo *int
	BlackElo *int

	Result string
	Date   string
	Event  string
	Site   string
	Round  string

	ECO     string
	Opening string

	PlyCount int

	Source GameSource
	Status GameStatus

	Partial bool

	StartFEN string
	PGN      string
	Hash     string
}

type UpdateGameParams struct {
	ID int64

	White    string
	Black    string
	WhiteElo *int
	BlackElo *int

	Result string
	Date   string
	Event  string
	Site   string
	Round  string

	ECO     string
	Opening string

	PlyCount int

	Source GameSource
	Status GameStatus

	Partial bool

	StartFEN string
	PGN      string
	Hash     string
}

type GameListFilter struct {
	Query string

	Player string
	Result string
	Source GameSource
	Status GameStatus

	Analysed *bool
	Partial  *bool

	DateFrom string
	DateTo   string

	Limit int

	Cursor string
}

type GameListPage struct {
	Games      []Game
	NextCursor string
}

type ImportGameResult struct {
	GameID  int64
	Created bool
}

type AnalysisRunStatus string

const (
	AnalysisRunCompleted AnalysisRunStatus = "completed"
	AnalysisRunCancelled AnalysisRunStatus = "cancelled"
	AnalysisRunFailed    AnalysisRunStatus = "failed"
)

// AnalysisRun represents a completed analysis of a game.
type AnalysisRun struct {
	ID int64

	GameID int64

	EngineID      string
	EngineName    string
	EngineVersion string
	SettingsJSON  string

	Status AnalysisRunStatus

	CreatedAt time.Time
}

// AnalysisPly represents analysis data for a single ply.
type AnalysisPly struct {
	RunID int64
	Ply   int

	CP   *int
	Mate *int

	BestUCI string
	BestSAN string

	Depth int

	Class    string
	WPLoss   float64
	Accuracy float64
}

type SaveAnalysisRunParams struct {
	GameID int64

	EngineID      string
	EngineName    string
	EngineVersion string
	SettingsJSON  string

	Plies []AnalysisPly
}

type EngineCache struct {
	EngineID string

	EngineName    string
	EngineVersion string
	OptionsJSON   string

	ProbedAt time.Time
}
