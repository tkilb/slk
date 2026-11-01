package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gammons/slk/internal/config"
)

func TestNewGrammarServiceRunsConfiguredCommand(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "fake-claude")
	if err := os.WriteFile(bin, []byte("#!/bin/sh\ncat >/dev/null; echo Fixed.\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := newGrammarService(config.Grammar{Command: bin}).Correct(context.Background(), "fixd")
	if err != nil || got != "Fixed." {
		t.Fatalf("got %q, %v", got, err)
	}
}
