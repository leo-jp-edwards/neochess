// Package analysis provides game analysis and evaluation services

package analysis

import (
	"context"
	"neochess/engines"
	"neochess/store"
	"neochess/uci"
	"sync"
)

// Service manages analysis jobs and live analysis
type Service struct {
	store interface {
		store.GameStore
		store.AnalysisStore
	}
	registry     *engines.Registry
	workerCount  int
	mu           sync.Mutex
	jobs         map[int64]*Job
	liveEngines  map[string]*liveAnalysis
	shutdownOnce sync.Once
}

// NewService creates a new analysis service
func NewService(s interface {
	store.GameStore
	store.AnalysisStore
}, r *engines.Registry) *Service {
	return &Service{
		store:       s,
		registry:    r,
		workerCount: 1,
		jobs:        make(map[int64]*Job),
		liveEngines: make(map[string]*liveAnalysis),
	}
}

// Job represents a running analysis job
type Job struct {
	ID       int64
	GameID   int64
	EngineID string
	Status   JobStatus
	Progress int
	Total    int
	cancel   context.CancelFunc // code smell here?
	resultCh chan JobProgress
}

// JobStatus indicates the state of an analysis job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusFailed    JobStatus = "failed"
)

// JobProgress reports analysis progress via SSE
type JobProgress struct {
	Done  int
	Total int
	Error error
}

// AnalysisOptions specifies the parameters for game analysis
type AnalysisOptions struct {
	EngineID string
	MoveTime int // milliseconds per position
	Depth    int // alternative to movetime
}

// StartAnalysis queues a game for full analysis
func (s *Service) StartAnalysis(ctx context.Context, gameID int64, opts AnalysisOptions) (*Job, error) {
	// TODO create job, add to queue, start worker if needed
	return nil, nil
}

// GetJob retrieves a job by ID
func (s *Service) GetJob(jobID int64) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[jobID]
	return j, ok
}

// CancelJob stops a running analysis job
func (s *Service) CancelJob(jobID int64) error {
	// TODO Call cancel func, wait for cleanup
	return nil
}

// Subscribe returns a channel for job progress updates
func (s *Service) Subscribe(jobID int64) (<-chan JobProgress, error) {
	// TODO Return channel that receives progress updates
	return nil, nil
}

// LiveAnalysisUpdate represents a single update from live analysis
type LiveAnalysisUpdate struct {
	Depth   int
	Score   int // centipawns, white POV
	Mate    int
	PV      []string // SAN moves
	MultiPV int
}

// StartLiveAnalysis begins live analysis of a position
// Cancels any previous live analysis for this engine
func (s *Service) StartLiveAnalysis(ctx context.Context, engineID, fen string, multiPV int) (<-chan LiveAnalysisUpdate, error) {
	// TODO get or start live engine, cancel previous search, start new search
	return nil, nil
}

// StopLiveAnalysis stops live analysis for an engine
func (s *Service) StopLiveAnalysis(engineID string) error {
	// TODO Stop current search
	return nil
}

// Shutdown gracefully stops all analysis jobs and engines
func (s *Service) Shutdown(ctx context.Context) {
	s.shutdownOnce.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		// Cancel all jobs
		for _, j := range s.jobs {
			if j.cancel != nil {
				j.cancel()
			}
		}

		// Close live engines
		for _, la := range s.liveEngines {
			if la.engine != nil {
			}
			la.engine.Close()
		}
	})
}

type liveAnalysis struct {
	engine *uci.Engine
	cancel context.CancelFunc
}
