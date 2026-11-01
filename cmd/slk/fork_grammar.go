// cmd/slk/fork_grammar.go
//
// Fork feature grammar-check: wires the claude-CLI proofreader into the
// TUI as its core.GrammarService.
package main

import (
	"time"

	"github.com/gammons/slk/internal/config"
	"github.com/gammons/slk/internal/core"
	"github.com/gammons/slk/internal/grammar"
)

func newGrammarService(cfg config.Grammar) core.GrammarService {
	return core.NewGrammarService(grammar.Corrector{
		Binary:  cfg.Command,
		Model:   cfg.Model,
		Prompt:  cfg.Prompt,
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}.Correct)
}
