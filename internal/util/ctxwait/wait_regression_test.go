package ctxwait

import (
	"context"
	"testing"
	"time"
)

// TestUntilReturnsOnCancel is a regression test ensuring Until stops waiting
// promptly when its context is cancelled, rather than sleeping for the full
// duration. Mirrors the production failure mode: a blocking time.Sleep that
// ignored ctx.Done().
func TestUntilReturnsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := Until(ctx, 5*time.Second)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if elapsed > 500*time.Millisecond {
		t.Fatalf("Until did not return promptly after cancel: elapsed=%s", elapsed)
	}
}

// TestUntilAlreadyCancelled verifies that a context already cancelled before
// the call returns immediately without waiting for the timer.
func TestUntilAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := Until(ctx, 5*time.Second)
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if elapsed > 100*time.Millisecond {
		t.Fatalf("Until did not return immediately for already-cancelled context: elapsed=%s", elapsed)
	}
}
