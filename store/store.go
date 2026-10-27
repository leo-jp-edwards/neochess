package store

import (
	"context"
)

type Store interface {
	GameStore
	AnalysisStore
	EngineCacheStore
	SettingsStore

	Close() error
	Ping(ctx context.Context) error
}

type GameStore interface {
	CreateGame(ctx context.Context, params CreateGameParams) (int64, error)
	InsertImportedGame(ctx context.Context, params CreateGameParams) (ImportGameResult, error)

	GetGame(ctx context.Context, id int64) (Game, error)
	GetGameIDByHash(ctx context.Context, hash string) (int64, error)

	ListGames(ctx context.Context, filter GameListFilter) (GameListPage, error)
	ListInProgressGames(ctx context.Context) ([]Game, error)

	UpdateGame(ctx context.Context, params UpdateGameParams) error
	DeleteGame(ctx context.Context, id int64) error
}

type AnalysisStore interface {
	SaveAnalysisRun(ctx context.Context, params SaveAnalysisRunParams) (int64, error)
	GetLatestAnalysisRun(ctx context.Context, gameID int64) (AnalysisRun, []AnalysisPly, error)
	ListAnalysisPlies(ctx context.Context, runID int64) ([]AnalysisPly, error)
	DeleteAnalysisForGame(ctx context.Context, gameID int64) error
}

type EngineCacheStore interface {
	UpsertEngineCache(ctx context.Context, cache EngineCache) error
	GetEngineCache(ctx context.Context, engineID string) (EngineCache, error)
	DeleteEngineCache(ctx context.Context, engineID string) error
}

type SettingsStore interface {
	SetSetting(ctx context.Context, key, value string) error
	GetSetting(ctx context.Context, key string) (string, error)
	DeleteSetting(ctx context.Context, key string) error
}
