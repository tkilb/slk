package ui

import (
	"context"
	"errors"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/core"
)

// grammarApp is an App in insert mode with the channel composer holding
// draft, whose proofreader records its inputs and answers with
// (corrected, err).
func grammarApp(t *testing.T, draft, corrected string, err error) (*App, *[]string) {
	t.Helper()
	a := newTestApp(t, append(insertOpts(), withActiveTeam("T1"), withMode(ModeInsert))...)
	a.focusedPanel = PanelMessages
	a.compose.SetDraftContext("T1", "C1", "")
	a.compose.SetValue(draft)
	calls := &[]string{}
	a.SetGrammarService(core.NewGrammarService(func(_ context.Context, text string) (string, error) {
		*calls = append(*calls, text)
		return corrected, err
	}))
	return a, calls
}

func pressCtrlG(a *App) tea.Cmd {
	_, cmd := a.Update(tea.KeyPressMsg{Code: 'g', Mod: tea.ModCtrl})
	return cmd
}

func TestGrammar_CtrlGCorrectsDraftWithToasts(t *testing.T) {
	a, calls := grammarApp(t, "ths is a tset", "This is a test.", nil)

	cmd := pressCtrlG(a)
	if got := statusbarText(a); !strings.Contains(got, grammarWorkingToast) {
		t.Fatalf("status bar = %q, want %q while running", got, grammarWorkingToast)
	}
	done := deliveredMsgs[grammarDoneMsg](t, cmd)
	if len(done) != 1 || len(*calls) != 1 || (*calls)[0] != "ths is a tset" {
		t.Fatalf("done=%v calls=%v", done, *calls)
	}

	a.Update(done[0])
	if got := a.compose.Value(); got != "This is a test." {
		t.Errorf("compose = %q, want corrected text", got)
	}
	if got := statusbarText(a); !strings.Contains(got, grammarDoneToast) {
		t.Errorf("status bar = %q, want %q", got, grammarDoneToast)
	}
	if a.mode != ModeInsert {
		t.Errorf("mode = %v, want to stay in insert mode", a.mode)
	}
}

func TestGrammar_ThreadComposeWhenThreadFocused(t *testing.T) {
	a, calls := grammarApp(t, "channel draft", "x", nil)
	a.threadVisible = true
	a.focusedPanel = PanelThread
	a.threadCompose.SetValue("thred reply")

	done := deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))
	if len(*calls) != 1 || (*calls)[0] != "thred reply" {
		t.Fatalf("calls = %v, want the thread draft", *calls)
	}
	a.Update(done[0])
	if a.threadCompose.Value() != "x" || a.compose.Value() != "channel draft" {
		t.Errorf("thread=%q channel=%q", a.threadCompose.Value(), a.compose.Value())
	}
}

func TestGrammar_StaleResultIsDiscarded(t *testing.T) {
	a, _ := grammarApp(t, "ths is", "This is.", nil)
	done := deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))

	a.compose.SetValue("ths is, and I kept typing")
	a.Update(done[0])

	if got := a.compose.Value(); got != "ths is, and I kept typing" {
		t.Errorf("compose = %q, newer text must be kept", got)
	}
	if got := statusbarText(a); !strings.Contains(got, "discarded") {
		t.Errorf("status bar = %q, want a discarded notice", got)
	}
}

func TestGrammar_FailureLeavesDraftAndToasts(t *testing.T) {
	a, _ := grammarApp(t, "ths is", "", errors.New("claude CLI not found on PATH"))
	done := deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))
	a.Update(done[0])

	if a.compose.Value() != "ths is" {
		t.Errorf("compose = %q, want unchanged", a.compose.Value())
	}
	if got := statusbarText(a); !strings.Contains(got, "Grammar check failed: claude CLI not found") {
		t.Errorf("status bar = %q", got)
	}
}

func TestGrammar_UnchangedText(t *testing.T) {
	a, _ := grammarApp(t, "All good.", "All good.", nil)
	done := deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))
	a.Update(done[0])
	if got := statusbarText(a); !strings.Contains(got, grammarUnchangedToast) {
		t.Errorf("status bar = %q, want %q", got, grammarUnchangedToast)
	}
}

func TestGrammar_EmptyDraftDoesNotCallService(t *testing.T) {
	a, calls := grammarApp(t, "  \n ", "x", nil)
	pressCtrlG(a)
	if len(*calls) != 0 {
		t.Errorf("service called %d times for a blank draft", len(*calls))
	}
	if got := statusbarText(a); !strings.Contains(got, grammarEmptyToast) {
		t.Errorf("status bar = %q", got)
	}
}

func TestGrammar_SecondPressWhileRunningIsRefused(t *testing.T) {
	a, calls := grammarApp(t, "ths is", "This is.", nil)
	first := pressCtrlG(a)
	if second := pressCtrlG(a); second == nil {
		t.Fatal("refusal should still clear its toast")
	}
	deliveredMsgs[grammarDoneMsg](t, first)
	if len(*calls) != 1 {
		t.Errorf("service called %d times, want 1", len(*calls))
	}
	if got := statusbarText(a); !strings.Contains(got, grammarBusyToast) {
		t.Errorf("status bar = %q, want %q", got, grammarBusyToast)
	}
}

func TestGrammar_PendingClearsSoItCanRunAgain(t *testing.T) {
	a, calls := grammarApp(t, "ths is", "This is.", nil)
	a.Update(deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))[0])
	a.compose.SetValue("anothr one")
	deliveredMsgs[grammarDoneMsg](t, pressCtrlG(a))
	if len(*calls) != 2 {
		t.Errorf("service called %d times, want 2", len(*calls))
	}
}

func TestGrammar_UnwiredServiceToasts(t *testing.T) {
	a, _ := grammarApp(t, "ths is", "x", nil)
	a.SetGrammarService(nil)
	pressCtrlG(a)
	if got := statusbarText(a); !strings.Contains(got, "not configured") {
		t.Errorf("status bar = %q", got)
	}
}

func TestGrammar_UploadInFlightBlocksIt(t *testing.T) {
	a, calls := grammarApp(t, "ths is", "x", nil)
	a.compose.SetUploading(true)
	pressCtrlG(a)
	if len(*calls) != 0 {
		t.Error("Ctrl+G must not run while an upload is in flight")
	}
}
