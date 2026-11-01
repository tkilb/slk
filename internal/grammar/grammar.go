// Package grammar proofreads text by piping it through the `claude` CLI
// (`claude -p`, no tools, no session). cmd/slk wires Corrector.Correct
// into the TUI as its core.GrammarService. Using the CLI reuses the
// user's existing Claude login, so slk never holds an API key.
package grammar

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// DefaultModel is a small, fast model: proofreading needs no more.
	DefaultModel = "haiku"
	// DefaultTimeout bounds one request.
	DefaultTimeout = 30 * time.Second
	// DefaultBinary is looked up on PATH.
	DefaultBinary = "claude"
)

// DefaultPrompt is the proofreader's system prompt.
const DefaultPrompt = "You are a proofreader for a chat message. " +
	"Fix spelling, grammar, and punctuation in the text the user sends. " +
	"Keep the original meaning, tone, wording, line breaks, markdown, code, " +
	"links, @mentions, #channels and :emoji: codes untouched unless a fix requires changing them. " +
	"Reply with ONLY the corrected text: no preamble, no explanation, no code fences. " +
	"If the text is already correct, return it unchanged."

// Corrector runs the claude CLI. The zero value uses the defaults.
type Corrector struct {
	Binary  string
	Model   string
	Prompt  string
	Timeout time.Duration
}

// ErrNoOutput means the CLI exited cleanly but printed nothing.
var ErrNoOutput = errors.New("claude returned no output")

// Correct returns text with spelling, grammar and punctuation fixed. The
// text goes over stdin, never argv, so it cannot be read as a flag or
// show up in the process list.
func (c Corrector) Correct(ctx context.Context, text string) (string, error) {
	bin := orDefault(c.Binary, DefaultBinary)
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin,
		"-p",
		"--model", orDefault(c.Model, DefaultModel),
		"--tools", "",
		"--no-session-persistence",
		"--system-prompt", orDefault(c.Prompt, DefaultPrompt),
	)
	// Run outside any repository so the CLI does not load project context.
	cmd.Dir = os.TempDir()
	cmd.Stdin = strings.NewReader(text)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if err := cmd.Run(); err != nil {
		switch {
		case errors.Is(err, exec.ErrNotFound):
			return "", fmt.Errorf("%s CLI not found on PATH", bin)
		case errors.Is(ctx.Err(), context.DeadlineExceeded):
			return "", fmt.Errorf("timed out after %s", timeout)
		case ctx.Err() != nil:
			return "", ctx.Err()
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	out := strings.TrimSpace(stdout.String())
	if out == "" {
		return "", ErrNoOutput
	}
	return out, nil
}

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
