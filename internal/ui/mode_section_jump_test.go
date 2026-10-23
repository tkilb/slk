package ui

import (
	"testing"

	"github.com/gammons/slk/internal/ui/sidebar"
)

func sectionJumpApp(t *testing.T) *App {
	t.Helper()
	a := newTestApp(t, withSize(120, 40), withChannels(
		sidebar.ChannelItem{ID: "C1", Name: "general", Type: "channel"},
		sidebar.ChannelItem{ID: "D1", Name: "alice", Type: "dm"},
	))
	a.SetSidebarCollapseWhenUnfocused(true)
	a.focusedPanel = PanelSidebar
	a.View()
	return a
}

func TestBraceKeysJumpSidebarSections(t *testing.T) {
	a := sectionJumpApp(t)
	a.sidebar.SelectByID("D1")

	a.Update(keyPress('}'))
	if name, ok := a.sidebar.IsSectionHeaderSelected(); !ok || name != "Channels" {
		t.Fatalf("} selected header = %q ok=%v, want Channels", name, ok)
	}
	a.Update(keyPress('{'))
	if name, ok := a.sidebar.IsSectionHeaderSelected(); !ok || name != "Direct Messages" {
		t.Errorf("{ selected header = %q ok=%v, want Direct Messages", name, ok)
	}
}

func TestBraceKeysAreIgnoredOutsideTheSidebar(t *testing.T) {
	a := sectionJumpApp(t)
	a.focusedPanel = PanelMessages
	a.View()
	before := sidebarRowKey(a)
	a.Update(keyPress('}'))
	if sidebarRowKey(a) != before {
		t.Errorf("} must not move the sidebar cursor when another panel has focus")
	}
}
