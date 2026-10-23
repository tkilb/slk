package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/ids"
)

func navHistoryApp() *App {
	app := NewApp()
	app.activeTeamID = "T1"
	app.setChannelLookupFuncForTest(func(channelID ids.ChannelID) (string, string, bool) {
		return string(channelID), "channel", true
	})
	app.Update(ChannelSelectedMsg{ID: "C1", Name: "a", Type: "channel"})
	app.Update(ChannelSelectedMsg{ID: "C2", Name: "b", Type: "channel"})
	return app
}

func historyTarget(t *testing.T, cmd tea.Cmd) string {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected a navigation command")
	}
	cs, ok := cmd().(ChannelSelectedMsg)
	if !ok || !cs.FromHistory {
		t.Fatalf("want a FromHistory ChannelSelectedMsg, got %T %+v", cmd(), cs)
	}
	return cs.ID
}

func TestCtrlONavigatesBackAndCtrlIForward(t *testing.T) {
	app := navHistoryApp()
	if got := historyTarget(t, app.handleNormalMode(tea.KeyPressMsg{Code: 'o', Mod: tea.ModCtrl})); got != "C1" {
		t.Fatalf("ctrl+o went to %q, want C1", got)
	}
	// Replay the history result so the forward stack has an entry.
	app.Update(ChannelSelectedMsg{ID: "C1", Name: "a", Type: "channel", FromHistory: true})
	if got := historyTarget(t, app.handleNormalMode(tea.KeyPressMsg{Code: 'i', Mod: tea.ModCtrl})); got != "C2" {
		t.Errorf("ctrl+i went to %q, want C2", got)
	}
}

func TestOriginalNavKeysStillWork(t *testing.T) {
	app := navHistoryApp()
	if got := historyTarget(t, app.handleNormalMode(tea.KeyPressMsg{Code: 'h', Mod: tea.ModCtrl})); got != "C1" {
		t.Errorf("ctrl+h went to %q, want C1", got)
	}
}
