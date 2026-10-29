package chessx

import (
	"errors"
	"fmt"
	"io"

	chess "github.com/corentings/chess/v2"
)

// ParsePGN parses a single PGN game into a Position.
func ParsePGN(r io.Reader) (*Position, error) {
	opt, err := chess.PGN(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPGN, err)
	}

	return &Position{
		game: chess.NewGame(opt),
	}, nil
}

// Scanner streams games from a PGN reader.
type Scanner struct {
	scanner *chess.Scanner
	current *chess.Game
	err     error
}

// NewScanner returns a PGN scanner suitable for multi-game PGN files.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		scanner: chess.NewScanner(r),
	}
}

// Scan advances to the next game.
// Keeps a bufio.Scanner-like interface for callers
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	s.current = nil

	if !s.scanner.HasNext() {
		return false
	}

	game, err := s.scanner.ParseNext()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return false
		}

		s.err = fmt.Errorf("%w: %v", ErrInvalidPGN, err)
		return false
	}

	s.current = game
	return true
}

// Position returns the current scanned game.
func (s *Scanner) Position() (*Position, error) {
	if s.current == nil {
		return nil, fmt.Errorf("%w: no current game", ErrInvalidPGN)
	}

	return &Position{
		game: s.current,
	}, nil
}

// Err returns the first scanner error, if any.
func (s *Scanner) Err() error {
	return s.err
}
