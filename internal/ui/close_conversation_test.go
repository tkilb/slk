package ui

import (
	"errors"
	"strings"
	"testing"

	"github.com/gammons/slk/internal/core"
	"github.com/gammons/slk/internal/ids"
	"github.com/gammons/slk/internal/ui/sidebar"
)

func closeDMApp(t *testing.T, closeFn func(ids.ChannelID) core.Cmd) *App {
	t.Helper()
	a := newTestApp(t, withSize(120, 40), withActiveTeam("T1"), withChannels(
		sidebar.ChannelItem{ID: "C1", Name: "general", Type: "channel"},
		sidebar.ChannelItem{ID: "D1", Name: "alice", Type: "dm"},
		sidebar.ChannelItem{ID: "G1", Name: "alice, bob", Type: "group_dm"},
	), withChannelService(core.ChannelServiceFuncs{CloseConversation: closeFn}))
	a.SetSidebarCollapseWhenUnfocused(true)
	a.focusedPanel = PanelSidebar
	a.View()
	return a
}

func sidebarHas(a *App, id string) bool {
	for _, it := range a.sidebar.Items() {
		if it.ID == id {
			return true
		}
	}
	return false
}

func TestXOnADMClosesIt(t *testing.T) {
	var closed []ids.ChannelID
	a := closeDMApp(t, func(id ids.ChannelID) core.Cmd {
		closed = append(closed, id)
		return func() core.Msg { return ConversationClosedResultMsg{ChannelID: string(id)} }
	})
	a.sidebar.SelectByID("D1")

	pressKey(t, a, keyPress('x')) // also feeds the result back through Update
	if len(closed) != 1 || closed[0] != "D1" {
		t.Fatalf("closed = %v, want [D1]", closed)
	}
	if sidebarHas(a, "D1") {
		t.Errorf("D1 should be removed once the close succeeded")
	}
	if got := statusbarText(a); !strings.Contains(got, "Closed alice") {
		t.Errorf("toast = %q, want it to say Closed alice", got)
	}
}

func TestXOnAGroupDMCloses(t *testing.T) {
	var closed []ids.ChannelID
	a := closeDMApp(t, func(id ids.ChannelID) core.Cmd {
		closed = append(closed, id)
		return nil
	})
	a.sidebar.SelectByID("G1")
	a.Update(keyPress('x'))
	if len(closed) != 1 || closed[0] != "G1" {
		t.Errorf("closed = %v, want [G1]", closed)
	}
}

func TestXOnAChannelRefusesWithAToast(t *testing.T) {
	called := false
	a := closeDMApp(t, func(ids.ChannelID) core.Cmd { called = true; return nil })
	a.sidebar.SelectByID("C1")
	a.Update(keyPress('x'))
	if called {
		t.Errorf("a regular channel must not be closed")
	}
	if got := statusbarText(a); !strings.Contains(got, "Only DMs") {
		t.Errorf("toast = %q, want an explanation", got)
	}
}

func TestXOutsideTheSidebarDoesNothing(t *testing.T) {
	called := false
	a := closeDMApp(t, func(ids.ChannelID) core.Cmd { called = true; return nil })
	a.sidebar.SelectByID("D1")
	a.focusedPanel = PanelMessages
	a.Update(keyPress('x'))
	if called {
		t.Errorf("x must only close from the sidebar")
	}
}

func TestCloseFailureKeepsTheRowAndToasts(t *testing.T) {
	a := closeDMApp(t, func(id ids.ChannelID) core.Cmd {
		return func() core.Msg { return ConversationClosedResultMsg{ChannelID: string(id), Err: errors.New("boom")} }
	})
	a.sidebar.SelectByID("D1")
	pressKey(t, a, keyPress('x'))
	if !sidebarHas(a, "D1") {
		t.Errorf("a failed close must keep the row")
	}
	if got := statusbarText(a); !strings.Contains(got, "boom") {
		t.Errorf("toast = %q, want the error", got)
	}
}

func TestWebSocketCloseRemovesTheRow(t *testing.T) {
	a := closeDMApp(t, nil)
	a.Update(ConversationClosedMsg{TeamID: "T1", ChannelID: "D1"})
	if sidebarHas(a, "D1") {
		t.Errorf("an im_close for the active workspace should remove D1")
	}
	a.Update(ConversationClosedMsg{TeamID: "OTHER", ChannelID: "G1"})
	if !sidebarHas(a, "G1") {
		t.Errorf("a close for another workspace must not touch this sidebar")
	}
}

func TestSlashCloseClosesTheCurrentDM(t *testing.T) {
	var closed []ids.ChannelID
	a, _ := slashApp(t, nil)
	a.SetChannelService(core.NewChannelService(core.ChannelServiceFuncs{
		CloseConversation: func(id ids.ChannelID) core.Cmd {
			closed = append(closed, id)
			return nil
		},
	}))
	a.sidebar.SetItems([]sidebar.ChannelItem{{ID: "C1", Name: "alice", Type: "dm"}})
	a.compose.SetValue("/close")
	pressEnter(a)
	if len(closed) != 1 || closed[0] != "C1" {
		t.Errorf("closed = %v, want [C1]", closed)
	}
}

func TestSlashCloseRefusesOnARegularChannel(t *testing.T) {
	called := false
	a, _ := slashApp(t, nil)
	a.SetChannelService(core.NewChannelService(core.ChannelServiceFuncs{
		CloseConversation: func(ids.ChannelID) core.Cmd { called = true; return nil },
	}))
	a.sidebar.SetItems([]sidebar.ChannelItem{{ID: "C1", Name: "general", Type: "channel"}})
	a.compose.SetValue("/close")
	pressEnter(a)
	if called {
		t.Errorf("/close must not close a regular channel")
	}
	if got := a.compose.Value(); got != "/close" {
		t.Errorf("a refused command should give its text back, composer = %q", got)
	}
}
