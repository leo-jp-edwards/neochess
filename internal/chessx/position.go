package chessx

import (
	"fmt"
	"strings"

	chess "github.com/corentings/chess/v2"
)

// FEN returns the current position in Forsyth-Edwards Notation.
func (p *Position) FEN() string {
	return p.game.Position().String()
}

// Turn returns the side to move.
func (p *Position) Turn() Color {
	if p.game.Position().Turn() == chess.White {
		return White
	}
	return Black
}

// LegalMoves returns all currently legal moves.
func (p *Position) LegalMoves() []Move {
	valid := p.game.ValidMoves()
	moves := make([]Move, 0, len(valid))

	for i := range valid {
		move := valid[i]
		moves = append(moves, Move{
			UCI: move.String(),
			SAN: p.sanForMove(&move),
		})
	}

	return moves
}

// IsLegal reports whether uci is legal in the current position.
func (p *Position) IsLegal(uci string) bool {
	uci = strings.TrimSpace(uci)
	if uci == "" {
		return false
	}

	for _, move := range p.game.ValidMoves() {
		if move.String() == uci {
			return true
		}
	}

	return false
}

// PushUCI applies a legal UCI move to the position.
func (p *Position) PushUCI(uci string) (Move, error) {
	uci = strings.TrimSpace(uci)
	if uci == "" {
		return Move{}, fmt.Errorf("%w: empty move", ErrInvalidMove)
	}

	for _, candidate := range p.game.ValidMoves() {
		if candidate.String() != uci {
			continue
		}

		move := candidate
		san := p.sanForMove(&move)
		if err := p.game.Move(&move, nil); err != nil {
			return Move{}, fmt.Errorf("%w: %v", ErrInvalidMove, err)
		}

		return Move{
			UCI: uci,
			SAN: san,
		}, nil
	}

	return Move{}, fmt.Errorf("%w: %q", ErrInvalidMove, uci)
}

// PushSAN applies a SAN move to the position.
func (p *Position) PushSAN(san string) (Move, error) {
	san = strings.TrimSpace(san)
	if san == "" {
		return Move{}, fmt.Errorf("%w: empty SAN move", ErrInvalidMove)
	}

	before := len(p.game.Moves())

	if err := p.game.PushNotationMove(san, chess.AlgebraicNotation{}, nil); err != nil {
		return Move{}, fmt.Errorf("%w: %v", ErrInvalidMove, err)
	}

	moves := p.game.Moves()
	if len(moves) != before+1 {
		return Move{}, fmt.Errorf("%w: move was not recorded", ErrInvalidMove)
	}

	applied := moves[len(moves)-1]

	return Move{
		UCI: applied.String(),
		SAN: san,
	}, nil
}

// PGN returns the current game as PGN text.
func (p *Position) PGN() string {
	return p.game.String()
}

// Outcome returns the current game outcome.
func (p *Position) Outcome() Outcome {
	switch p.game.Outcome() {
	case chess.WhiteWon:
		return WhiteWon
	case chess.BlackWon:
		return BlackWon
	case chess.Draw:
		return Draw
	default:
		return NoOutcome
	}
}

// Method returns the current game-end method.
func (p *Position) Method() Method {
	switch p.game.Method() {
	case chess.Checkmate:
		return MethodCheckmate
	case chess.Stalemate:
		return MethodStalemate
	case chess.InsufficientMaterial:
		return MethodInsufficientMat
	case chess.ThreefoldRepetition:
		return MethodThreefold
	case chess.FivefoldRepetition:
		return MethodFivefold
	case chess.FiftyMoveRule:
		return MethodFiftyMoveRule
	case chess.SeventyFiveMoveRule:
		return MethodSeventyFiveMove
	default:
		if p.Outcome() == NoOutcome {
			return MethodNone
		}
		return MethodOther
	}
}

// Moves returns the played move list.
func (p *Position) Moves() []Move {
	raw := p.game.Moves()
	out := make([]Move, 0, len(raw))

	replay := chess.NewGame()

	for _, move := range raw {
		san := notationSAN(replay, move)

		out = append(out, Move{
			UCI: move.String(),
			SAN: san,
		})

		_ = replay.Move(move, nil)
	}

	return out
}

func (p *Position) sanForMove(move *chess.Move) string {
	return notationSAN(p.game, move)
}

func notationSAN(game *chess.Game, move *chess.Move) string {
	notation := chess.AlgebraicNotation{}
	return notation.Encode(game.Position(), move)
}
