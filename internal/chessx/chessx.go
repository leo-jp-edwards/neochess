// Package chessx provides an adapter over the chess library, isolating
// the rest of the application from the specific library implementation

package chessx

// Color represents a player color
type Color int

const (
	White Color = iota
	Black
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
type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

// Position represents a chess position
type Position struct {
	// Underlying library position
}

// NewPosition creates the standard starting position.
func NewPosition() *Position {
	// TODO: Wrap library's starting position
	return &Position{}
}

// FromFEN creates a position from FEN notation.
func FromFEN(fen string) (*Position, error) {
	// TODO: Parse FEN using library
	return nil, nil
}

// FEN returns the FEN representation of the position.
func (p *Position) FEN() string {
	// TODO: Implement
	return ""
}

// Turn returns whose turn it is to move.
func (p *Position) Turn() Color {
	// TODO: Implement
	return White
}

// LegalMoves returns all legal moves in the position.
func (p *Position) LegalMoves() []Move {
	// TODO: Implement
	return nil
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

// Outcome returns the game outcome if the game is over.
func (p *Position) Outcome() Outcome {
	// TODO: Implement
	return OutcomeOngoing
}

// Outcome represents the result of a game.
type Outcome int

const (
	OutcomeOngoing Outcome = iota
	OutcomeWhiteWins
	OutcomeBlackWins
	OutcomeDraw
)

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
