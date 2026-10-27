// Package pgnio handles PGN import/export with tolerance for real-world messy files

package pgnio

import (
	"context"
	"io"
	"neochess/internal/chessx"
)

// ImportedGame represents a game parsed from PGN
type ImportedGame struct {
	Headers  map[string]string
	Root     *chessx.GameNode
	Mainline []chessx.Move
	StartFEN string
	PGN      string // Original PGN text
	Hash     string // Deduplication hash
	Partial  bool   // True if truncated due to illegal move
	Error    string // Non-empty if parsing had issues
}

// ImportResult summarises an import operation
type ImportResult struct {
	Added      int
	Duplicates int
	Skipped    int
	Partial    int
	Errors     []string
}

// Import parses PGN from a reader, streaming games one at a time
func Import(ctx context.Context, r io.Reader) (<-chan *ImportedGame, error) {
	// TODO Stream parse PGN, yielding games one by one
	// Tolerant parsing: ignore unknown tags, accept missing roster fields
	// truncate at illegal moves
	return nil, nil
}

// Export writes games as PGN to the writer
func Export(w io.Writer, games []*ImportedGame) error {
	// TODO write standard PGN format
	return nil
}

// ComputeHash generates a deduplication hash for a game
// Based on normalised headers + mainline moves + starting FEN
func ComputeHash(headers map[string]string, moves []chessx.Move, startFEN string) string {
	// TODO implement deterministic hash
	return ""
}

// NormaliseHeaders cleans up header values for consistent hashing
func NormaliseHeaders(headers map[string]string) map[string]string {
	// TODO trim whitespace, normalise date formate etc
	return headers
}
