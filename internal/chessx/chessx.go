// Package chessx provides an adapter over the chess library, isolating
// the rest of the application from the specific library implementation

package chessx

import (
	"fmt"
	"strings"

	chess "github.com/corentings/chess/v2"
)

// Color represents a player color
type Color string

const (
	White Color = "white"
	Black Color = "black"
)

// Outcome identifies the game result from the current position
type Outcome string

const (
	NoOutcome Outcome = "*"
	WhiteWon  Outcome = "1-0"
	BlackWon  Outcome = "0-1"
	Draw      Outcome = "1/2-1/2"
)

// Method describes how the game ended.
type Method string

const (
	MethodNone            Method = ""
	MethodCheckmate       Method = "checkmate"
	MethodStalemate       Method = "stalemate"
	MethodInsufficientMat Method = "insufficient_material"
	MethodThreefold       Method = "threefold_repetition"
	MethodFivefold        Method = "fivefold_repetition"
	MethodFiftyMoveRule   Method = "fifty_move_rule"
	MethodSeventyFiveMove Method = "seventy_five_move_rule"
	MethodOther           Method = "other"
)

// PieceType represents a chess piece type
type PieceType int

const (
	NoPiece PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// Square represents a board square (0-63, a1=0, h8=63)
type Square int

// Move represents a chess move
// UCI is the canonical internal representation used for storage and UCI engines,
// for example "e2e4", "e7e8q"
type Move struct {
	UCI string
	SAN string
}

// Position represents a chess position
// Intentionally do not expose github.com/corentings/chess/v2 types
type Position struct {
	game *chess.Game
}

// NewPosition creates the standard starting position.
func NewPosition() *Position {
	return &Position{
		game: chess.NewGame(),
	}
}

// FromFEN creates a position from FEN notation.
func FromFEN(fen string) (*Position, error) {
	fen = strings.TrimSpace(fen)
	if fen == "" {
		return nil, fmt.Errorf("%w: empty FEN", ErrInvalidFEN)
	}

	opt, err := chess.FEN(fen)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidFEN, err)
	}

	return &Position{game: chess.NewGame(opt)}, nil
}

// Clone returns an independent copy of the position
// This uses FEN as the stable boundary instead of sharing upstream mutable state
func (p *Position) Clone() (*Position, error) {
	if p == nil || p.game == nil {
		return NewPosition(), nil
	}

	return FromFEN(p.FEN())
}

// LegalMovesFrom returns legal moves from a specific square.
func (p *Position) LegalMovesFrom(sq Square) []Move {
	// TODO: Implement
	return nil
}

// MakeMove applies a move and returns the new position.
func (p *Position) MakeMove(m Move) (*Position, error) {
	// TODO: Implement
	return nil, nil
}

// IsCheck returns true if the side to move is in check.
func (p *Position) IsCheck() bool {
	// TODO: Implement
	return false
}

// MoveToUCI converts a move to UCI notation.
func (m Move) ToUCI() string {
	// TODO: Implement (e.g., "e2e4", "e7e8q")
	return ""
}

// MoveFromUCI parses a move from UCI notation.
func MoveFromUCI(s string) (Move, error) {
	// TODO: Implement
	return Move{}, nil
}

// MoveToSAN converts a move to SAN notation given the position.
func (p *Position) MoveToSAN(m Move) string {
	// TODO: Implement
	return ""
}

// MoveFromSAN parses a move from SAN notation.
func (p *Position) MoveFromSAN(s string) (Move, error) {
	// TODO: Implement
	return Move{}, nil
}

// GameNode represents a node in a move tree (for variations/studies).
type GameNode struct {
	Move        *Move
	Comment     string
	NAGs        []int
	Annotations []string // [%cal], [%csl] annotations
	Children    []*GameNode
	Parent      *GameNode
}

// Mainline returns the mainline moves from this node.
func (n *GameNode) Mainline() []Move {
	// TODO: Implement
	return nil
}
