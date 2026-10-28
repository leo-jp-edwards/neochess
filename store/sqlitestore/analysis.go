package sqlitestore

import (
	"context"
	"neochess/store"
)

func (s *Store) SaveAnalysisRun(ctx context.Context, params store.SaveAnalysisRunParams) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetLatestAnalysisRun(ctx context.Context, gameID int64) (store.AnalysisRun, []store.AnalysisPly, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) ListAnalysisPlies(ctx context.Context, runID int64) ([]store.AnalysisPly, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) DeleteAnalysisForGame(ctx context.Context, gameID int64) error {
	//TODO implement me
	panic("implement me")
}
