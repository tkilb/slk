package grammar

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// fakeCLI writes an executable shell script that stands in for claude.
func fakeCLI(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "claude")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCorrectReturnsTrimmedStdout(t *testing.T) {
	bin := fakeCLI(t, `cat >/dev/null; printf 'Fixed text.\n'`)
	got, err := Corrector{Binary: bin}.Correct(context.Background(), "fixd text")
	if err != nil || got != "Fixed text." {
		t.Fatalf("got %q, %v", got, err)
	}
}

func TestCorrectSendsTextOnStdinAndFixedArgs(t *testing.T) {
	argsFile := filepath.Join(t.TempDir(), "args")
	bin := fakeCLI(t, `printf '%s\n' "$@" > `+argsFile+`; cat`)
	text := "--dangerous -p text"
	got, err := Corrector{Binary: bin, Model: "m1", Prompt: "P"}.Correct(context.Background(), text)
	if err != nil || got != text {
		t.Fatalf("stdin not echoed: %q, %v", got, err)
	}
	args, _ := os.ReadFile(argsFile)
	for _, want := range []string{"-p", "--model\nm1", "--tools\n", "--no-session-persistence", "--system-prompt\nP"} {
		if !strings.Contains(string(args), want) {
			t.Errorf("args missing %q:\n%s", want, args)
		}
	}
	if strings.Contains(string(args), "dangerous") {
		t.Errorf("user text leaked into argv:\n%s", args)
	}
}

func TestCorrectErrors(t *testing.T) {
	t.Run("missing binary", func(t *testing.T) {
		_, err := Corrector{Binary: "slk-no-such-claude"}.Correct(context.Background(), "x")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("non-zero exit carries stderr", func(t *testing.T) {
		bin := fakeCLI(t, `echo 'not logged in' >&2; exit 1`)
		_, err := Corrector{Binary: bin}.Correct(context.Background(), "x")
		if err == nil || !strings.Contains(err.Error(), "not logged in") {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("empty output", func(t *testing.T) {
		bin := fakeCLI(t, `cat >/dev/null`)
		_, err := Corrector{Binary: bin}.Correct(context.Background(), "x")
		if !errors.Is(err, ErrNoOutput) {
			t.Fatalf("err = %v", err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		bin := fakeCLI(t, `exec sleep 5`)
		start := time.Now()
		_, err := Corrector{Binary: bin, Timeout: 100 * time.Millisecond}.Correct(context.Background(), "x")
		if err == nil || !strings.Contains(err.Error(), "timed out") {
			t.Fatalf("err = %v", err)
		}
		if time.Since(start) > 3*time.Second {
			t.Fatal("did not stop at the timeout")
		}
	})
}
