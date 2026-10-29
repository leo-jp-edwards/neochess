package chessx

import "errors"

var (
	ErrInvalidFEN  = errors.New("chessx: invalid FEN")
	ErrInvalidMove = errors.New("chessx: invalid move")
	ErrInvalidPGN  = errors.New("chessx: invalid PGN")
)
