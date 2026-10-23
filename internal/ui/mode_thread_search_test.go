package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func threadSearchApp(t *testing.T) *App {
	t.Helper()
	a := newTestApp(t, withSize(120, 40), withActiveChannel("C1"))
	focusThreadPanel(t, a)
	return a
}

func TestSlashOnThreadOpensThreadSearch(t *testing.T) {
	a := threadSearchApp(t)
	a.Update(keyPress('/'))
	if a.mode != ModeThreadSearch {
		t.Fatalf("mode = %v, want ModeThreadSearch", a.mode)
	}
}

func TestThreadSearchTypingSelectsAndReportsMatches(t *testing.T) {
	a := threadSearchApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "reply")
	if got := a.statusbar.Search(); got != "/reply  2/2" {
		t.Errorf("status = %q, want /reply  2/2 (cursor starts on the newest reply)", got)
	}
	if r := a.threadPanel.SelectedReply(); r == nil || r.TS != "12.0" {
		t.Errorf("selected = %+v, want reply-2", r)
	}
	typeSidebarFilter(a, "x")
	if got := a.statusbar.Search(); got != "/replyx  no matches" {
		t.Errorf("status = %q, want no matches", got)
	}
}

func TestThreadSearchEnterThenNStepsAndWraps(t *testing.T) {
	a := threadSearchApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "reply")
	a.Update(keyCode(tea.KeyEnter))
	if a.mode != ModeNormal || !a.threadPanel.SearchActive() {
		t.Fatalf("mode=%v active=%v, want normal with search kept", a.mode, a.threadPanel.SearchActive())
	}
	a.Update(keyPress('n')) // wraps from 2/2 to 1/2
	if r := a.threadPanel.SelectedReply(); r == nil || r.TS != "11.0" {
		t.Errorf("n should wrap to reply-1, got %+v", r)
	}
	if got := a.statusbar.Search(); got != "/reply  1/2" {
		t.Errorf("status = %q", got)
	}
	a.Update(keyPress('N')) // back to 2/2
	if r := a.threadPanel.SelectedReply(); r == nil || r.TS != "12.0" {
		t.Errorf("N should return to reply-2, got %+v", r)
	}
}

func TestThreadSearchEscClearsBeforeClosingThread(t *testing.T) {
	a := threadSearchApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "reply")
	a.Update(keyCode(tea.KeyEnter))

	a.Update(keyCode(tea.KeyEscape))
	if a.threadPanel.SearchActive() || a.statusbar.Search() != "" {
		t.Errorf("first Esc should clear the search and its status")
	}
	if !a.threadVisible {
		t.Fatalf("first Esc must not close the thread")
	}
	a.Update(keyCode(tea.KeyEscape))
	if a.threadVisible {
		t.Errorf("second Esc should close the thread as before")
	}
}

func TestThreadSearchEscInPromptClears(t *testing.T) {
	a := threadSearchApp(t)
	a.Update(keyPress('/'))
	typeSidebarFilter(a, "reply")
	a.Update(keyCode(tea.KeyEscape))
	if a.mode != ModeNormal || a.threadPanel.SearchActive() || a.statusbar.Search() != "" {
		t.Errorf("mode=%v active=%v status=%q", a.mode, a.threadPanel.SearchActive(), a.statusbar.Search())
	}
	if !a.threadVisible {
		t.Errorf("Esc in the prompt must not close the thread")
	}
}

func TestChannelSearchNKeysUnaffectedByThreadSearch(t *testing.T) {
	a := threadSearchApp(t)
	a.focusedPanel = PanelMessages
	a.Update(keyPress('n')) // no channel search, no thread search: a no-op
	if a.threadPanel.SearchActive() {
		t.Errorf("n on the messages pane must not touch the thread search")
	}
}
