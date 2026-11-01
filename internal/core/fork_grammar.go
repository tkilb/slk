// internal/core/fork_grammar.go
//
// Fork feature grammar-check: the port the TUI calls to proofread a
// compose draft. Kept out of ports.go/adapters.go so upstream edits to
// those files rebase cleanly.
package core

import (
	"context"
	"errors"
)

// GrammarService proofreads a compose draft.
type GrammarService interface {
	// Correct returns text with spelling, grammar and punctuation fixed.
	// It blocks, so callers run it off the Update goroutine.
	Correct(ctx context.Context, text string) (string, error)
}

// NewGrammarService builds a GrammarService from a closure. A nil
// correct fails with errors.ErrUnsupported.
func NewGrammarService(correct func(ctx context.Context, text string) (string, error)) GrammarService {
	return grammarAdapter{correct: correct}
}

type grammarAdapter struct {
	correct func(ctx context.Context, text string) (string, error)
}

func (g grammarAdapter) Correct(ctx context.Context, text string) (string, error) {
	if g.correct == nil {
		return "", errors.ErrUnsupported
	}
	return g.correct(ctx, text)
}
