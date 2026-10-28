package sqlitestore

import (
	"context"
	"neochess/store"
)

func (s *Store) UpsertEngineCache(ctx context.Context, cache store.EngineCache) error {
	//TODO implement me
	panic("implement me")
}

func (s *Store) GetEngineCache(ctx context.Context, engineID string) (store.EngineCache, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Store) DeleteEngineCache(ctx context.Context, engineID string) error {
	//TODO implement me
	panic("implement me")
}
