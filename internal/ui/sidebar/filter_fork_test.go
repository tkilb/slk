package sidebar

import "testing"

func filterModel() *Model {
	m := New([]ChannelItem{
		{ID: "C1", Name: "general", Type: "channel"},
		{ID: "C2", Name: "random", Type: "channel"},
		{ID: "D1", Name: "alice", Type: "dm"},
	})
	m.SetCollapseWhenUnfocused(true)
	m.SetFocused(true)
	return &m
}

func TestSelectFirstChannel(t *testing.T) {
	m := filterModel()
	m.SetFilter("ran")
	if !m.SelectFirstChannel() {
		t.Fatalf("expected a channel row to select")
	}
	if got := m.SelectedID(); got != "C2" {
		t.Errorf("SelectedID = %q, want C2", got)
	}
}

func TestSelectFirstChannel_NoMatches(t *testing.T) {
	m := filterModel()
	m.SetFilter("zzz")
	if m.SelectFirstChannel() {
		t.Errorf("no channel rows should report false")
	}
}

func TestClearFilter_KeepsSelection(t *testing.T) {
	m := filterModel()
	m.SetFilter("ran")
	m.SelectFirstChannel()
	m.ClearFilter()
	if m.Filter() != "" {
		t.Errorf("filter should be empty, got %q", m.Filter())
	}
	if got := m.SelectedID(); got != "C2" {
		t.Errorf("selection should stay on the match, got %q", got)
	}
	if n := len(m.VisibleItems()); n != 3 {
		t.Errorf("all items should be visible again, got %d", n)
	}
}

func TestFilterClearsWhenSidebarLosesFocus(t *testing.T) {
	m := filterModel()
	m.SetFilter("ran")
	m.SetFocused(false)
	if m.Filter() != "" {
		t.Errorf("blur should clear the filter, got %q", m.Filter())
	}
}
