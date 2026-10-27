// Package play manages chess game sessions against UCI engines
package play

import (
	"context"
	"neochess/engines"
	"neochess/internal/chessx"
	"neochess/store"
	"neochess/uci"
	"sync"
)

// SessionManager coordinates active play sessions
type SessionManager struct {
	store    *store.Store
	registry *engines.Registry
	mu       sync.Mutex
	sessions map[int64]*Session
}

// NewSessionManager creates a new session manager
func NewSessionManager(s *store.Store, r *engines.Registry) *SessionManager {
	return &SessionManager{
		store:    s,
		registry: r,
		sessions: make(map[int64]*Session),
	}
}

// Session represents an active game against an engine
type Session struct {
	GameID   int64
	Position *chessx.Position
	Moves    []chessx.Move
	Engine   *uci.Engine
	Status   SessionStatus
	Limits   uci.Limits
	mu       sync.Mutex
}

// SessionStatus indicates the current stat of a session
type SessionStatus string

const (
	StatusPlayerTurn   SessionStatus = "player_turn"
	StatusEngineTurn   SessionStatus = "engine_turn"
	StatusFinished     SessionStatus = "finished"
	StatusEngineFailed SessionStatus = "engine_failed"
)

// NewGameOptions specifies parameters for starting a new game
type NewGameOptions struct {
	EngineID   string
	PlayerSide chessx.Color // White, Black, Random
	Limits     uci.Limits
	EloLimit   int // 0 = no limit
}

// NewGame starts a new game session
func (m *SessionManager) NewGame(ctx context.Context, opts NewGameOptions) (*Session, error) {
	// TODO: Create game in store, start engine, create session
	return nil, nil
}

// ResumeGame resumes an in-progress game
func (m *SessionManager) ResumeGame(ctx context.Context, opts NewGameOptions) (*Session, error) {
	// TODO Load game from store, restore position, start engine
	return nil, nil
}

// GetSession retrieves an active session by game ID
func (m *SessionManager) GetSession(gameID int64) (*Session, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[gameID]
	return s, ok
}

// SelectSquare handles square selection, returning legal targets
func (s *Session) SelectSquare(square chessx.Square) []chessx.Square {
	// TODO return legal move destinations from this square
	return nil
}

// MakeMove attempts to make a move, triggering engine response if valid
func (s *Session) MakeMove(ctx context.Context, from, to chessx.Square, promo chessx.PieceType) error {
	// TODO validate move, update position, save to store, trigger engine turn
	return nil
}

// EngineMove returns a channel for the engine's response
func (s *Session) EngineMove(ctx context.Context) (<-chan EngineUpdate, error) {
	// TODO Start engine search, return channel with progress and final move
	return nil, nil
}

// EngineUpdate represents progress or result of engine thinking
type EngineUpdate struct {
	Thinking bool
	Info     *uci.Info
	Move     *chessx.Move
	GameOver bool
	Result   string
}

// Resign ends the game with a resignation
func (s *Session) Resign(ctx context.Context) error {
	// TODO Update game result in store, close engine
	return nil
}

// Shutdown gracefully closes all sessions
func (m *SessionManager) Shutdown(ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.sessions {
		if s.Engine != nil {
			s.Engine.Close()
		}
	}
	m.sessions = make(map[int64]*Session)
}
