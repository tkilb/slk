// internal/ui/fork_grammar.go
//
// Fork feature grammar-check: Ctrl+G in insert mode proofreads the
// compose draft through core.GrammarService and replaces it with the
// corrected text. The request runs off the Update goroutine; a toast
// reports "Grammar correcting…" and then "Corrected".
//
// The result is applied only if the draft still holds exactly the text
// that was sent. Anything typed (or a channel switch that swapped the
// draft) in the meantime wins: the correction is discarded, never
// merged over newer text.
package ui

import (
	"context"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/core"
)

const (
	grammarWorkingToast   = "Grammar correcting…"
	grammarDoneToast      = "Corrected"
	grammarUnchangedToast = "No corrections needed"
	grammarStaleToast     = "Draft changed while checking, correction discarded"
	grammarEmptyToast     = "Nothing to check"
	grammarBusyToast      = "Grammar check already running"

	grammarToastDuration      = 2 * time.Second
	grammarFailureToastLength = 4 * time.Second
)

// grammarDoneMsg reports a finished grammar check. Original is the draft
// that was sent, so the reducer can tell whether it is still current.
type grammarDoneMsg struct {
	Panel     Panel
	Original  string
	Corrected string
	Err       error
}

// SetGrammarService wires Ctrl+G's proofreader. nil disables it.
func (a *App) SetGrammarService(s core.GrammarService) {
	a.grammar = s
}

// grammarTarget returns the compose Ctrl+G acts on, like Ctrl+E and Ctrl+U.
func (a *App) grammarTarget() (Panel, composeDraft) {
	if a.focusedPanel == PanelThread && a.threadVisible {
		return PanelThread, &a.threadCompose
	}
	return PanelMessages, &a.compose
}

// composeDraft is the part of the compose model Ctrl+G uses.
type composeDraft interface {
	Value() string
	SetValue(string)
}

func (a *App) checkGrammar() tea.Cmd {
	if a.grammar == nil {
		return toastWithClear(a, "Grammar check is not configured", grammarFailureToastLength)
	}
	panel, target := a.grammarTarget()
	text := target.Value()
	if strings.TrimSpace(text) == "" {
		return toastWithClear(a, grammarEmptyToast, grammarToastDuration)
	}
	if a.grammarPending[panel] {
		return toastWithClear(a, grammarBusyToast, grammarToastDuration)
	}
	if a.grammarPending == nil {
		a.grammarPending = map[Panel]bool{}
	}
	a.grammarPending[panel] = true
	toastUntilReplaced(a, grammarWorkingToast)

	svc := a.grammar
	return func() tea.Msg {
		corrected, err := svc.Correct(context.Background(), text)
		return grammarDoneMsg{Panel: panel, Original: text, Corrected: corrected, Err: err}
	}
}

var reduceGrammar reducerFunc = func(a *App, msg tea.Msg) (tea.Cmd, bool) {
	m, ok := msg.(grammarDoneMsg)
	if !ok {
		return nil, false
	}
	delete(a.grammarPending, m.Panel)

	if m.Err != nil {
		return toastWithClear(a, "Grammar check failed: "+m.Err.Error(), grammarFailureToastLength), true
	}
	target := composeDraft(&a.compose)
	if m.Panel == PanelThread {
		target = &a.threadCompose
	}
	if target.Value() != m.Original {
		return toastWithClear(a, grammarStaleToast, grammarFailureToastLength), true
	}
	if m.Corrected == m.Original {
		return toastWithClear(a, grammarUnchangedToast, grammarToastDuration), true
	}
	target.SetValue(m.Corrected)
	return toastWithClear(a, grammarDoneToast, grammarToastDuration), true
}
