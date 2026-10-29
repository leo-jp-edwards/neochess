package chessx

import "testing"

func TestNewPosition(t *testing.T) {
	pos := NewPosition()

	if got := pos.Turn(); got != White {
		t.Fatalf("turn = %q, want %q", got, White)
	}

	if got := len(pos.LegalMoves()); got != 20 {
		t.Fatalf("legal moves = %d, want 20", got)
	}

	if got := pos.Outcome(); got != NoOutcome {
		t.Fatalf("outcome = %q, want %q", got, NoOutcome)
	}
}

func TestPushUCI(t *testing.T) {
	pos := NewPosition()

	move, err := pos.PushUCI("e2e4")
	if err != nil {
		t.Fatalf("PushUCI: %v", err)
	}

	if move.UCI != "e2e4" {
		t.Fatalf("move UCI = %q, want e2e4", move.UCI)
	}

	if move.SAN != "e4" {
		t.Fatalf("move SAN = %q, want e4", move.SAN)
	}

	if got := pos.Turn(); got != Black {
		t.Fatalf("turn = %q, want %q", got, Black)
	}
}

func TestIllegalMove(t *testing.T) {
	pos := NewPosition()

	if _, err := pos.PushUCI("e2e5"); err == nil {
		t.Fatal("PushUCI illegal move returned nil error")
	}
}

func TestFromFEN(t *testing.T) {
	pos, err := FromFEN("8/8/8/8/8/8/8/K6k w - - 0 1")
	if err != nil {
		t.Fatalf("FromFEN: %v", err)
	}

	if got := pos.Turn(); got != White {
		t.Fatalf("turn = %q, want %q", got, White)
	}
}
