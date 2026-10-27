// Package uci provides a UCI protocol client for chess engine communication
package uci

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// Sentinel errors for UCI ops
var (
	ErrEngineExited      = errors.New("uci: engine process exited unexpectedly")
	ErrTimeout           = errors.New("uci: operation timed out")
	ErrUnsupportedOption = errors.New("uci: unsupported option")
)

// Engine represents a running UCI engine process
type Engine struct {
	path    string
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	options map[string]Option
}

// Option describes a UCI option advertised by the engine
type Option struct {
	Name    string
	Type    string // "check", "spin", "combo", "button", "string"
	Default string
	Min     int
	Max     int
	Vars    []string // for combo type
}

// Limits specifies search constraints for the go command
type Limits struct {
	MoveTime int // milliseconds
	Depth    int
	Nodes    int64
}

// GoOpts specifies additional options for the go command
type GoOpts struct {
	MultiPV     int
	SearchMoves []string // UCI move strings to restrict search
}

// Info represents parsed info lines from the engine
type Info struct {
	Depth    int
	SelDepth int
	MultiPV  int
	Score    Score
	Nodes    int64
	NPS      int64
	Time     int // milliseconds
	PV       []string
	BestMove string // non-empty only for the final message
	Ponder   string
}

// Score represents an engine evaluation
type Score struct {
	CP         int
	Mate       int // positive = engine mates in N, negative = engine mated in N
	LowerBound bool
	UpperBound bool
}

// Position represents a chess position for the engine
type Position struct {
	FEN   string   // starting position FEN (empty for startpos)
	Moves []string // UCI moves from the position
}

// Start launches a UCI engine process and performs the initial handshake
func Start(ctx context.Context, path string, args []string) (*Engine, error) {
	// TODO start process, send "uci", parse id/option lines until "uciok"
	// Send "isready", wait for "readyok" with timeout
	return nil, fmt.Errorf("uci: not implemented")
}

// Options returns the options advertised by the engine during handshake
func (e *Engine) Options() map[string]Option {
	return e.options
}

// SetOption sends a setoption command to the engine
func (e *Engine) SetOption(name, value string) error {
	// TODO Validate option exists, send "setoption name X value Y"
	return nil
}

// NewGame sends ucinewgame and waits for readyok
func (e *Engine) NewGame(ctx context.Context) error {
	// TODO Send "ucinewgame" then "isready", wait for "readyok"
	return nil
}

// Go starts a search and returns a channel of info updates
// The channel is closed after BestMove is received or context is cancelled
func (e *Engine) Go(ctx context.Context, pos Position, Lim Limits, opts GoOpts) (<-chan Info, error) {
	// TODO Send "position" command, then "go" with limits
	// Spawn goroutine to parse info lines and send to channel
	// On context cancel, send "stop" and wait for bestmove
	return nil, nil
}

// Stop sends the stop command and waits for bestmove
func (e *Engine) Stop(ctx context.Context) error {
	// TODO Send "stop", wait for "bestmove" with grace timeout
	return nil
}

// Close terminates the engine process gracefully
func (e *Engine) Close() error {
	// TODO Send "quit", wait briefly, then kill if needed
	// Ensure all goroutines exit, idempotent
	return nil
}
