package sidebar

import "testing"

func TestRemoveItem_DropsTheRowAndReportsIt(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	if !m.RemoveItem("D1") {
		t.Fatalf("RemoveItem should report it removed D1")
	}
	for _, it := range m.Items() {
		if it.ID == "D1" {
			t.Errorf("D1 should be gone from Items")
		}
	}
	for _, it := range m.VisibleItems() {
		if it.ID == "D1" {
			t.Errorf("D1 should be gone from VisibleItems")
		}
	}
}

func TestRemoveItem_UnknownIDIsANoOp(t *testing.T) {
	m := unfocusedCollapseModel()
	if m.RemoveItem("nope") {
		t.Errorf("an unknown ID must report false")
	}
	if len(m.Items()) != 4 {
		t.Errorf("items changed: %d", len(m.Items()))
	}
}

func TestRemoveItem_CursorLandsOnTheNextRow(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("D1") // nav: ... Direct Messages, D1, D2, Channels ...
	m.RemoveItem("D1")
	if got := m.SelectedID(); got != "D2" {
		t.Errorf("cursor should move to the next row (D2), got %q", got)
	}
}

func TestRemoveItem_LastRowOfTheListClampsTheCursor(t *testing.T) {
	m := New([]ChannelItem{{ID: "D1", Name: "alice", Type: "dm"}})
	m.SetCollapseWhenUnfocused(true)
	m.SetFocused(true)
	m.SelectByID("D1")
	m.RemoveItem("D1")
	if m.cursor < 0 || m.cursor >= len(m.nav) {
		t.Errorf("cursor %d out of range for %d nav rows", m.cursor, len(m.nav))
	}
}

func TestRemoveItem_LeavesOtherCursorsAlone(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("C1")
	m.RemoveItem("D1")
	if got := m.SelectedID(); got != "C1" {
		t.Errorf("cursor on another row must stay put, got %q", got)
	}
}
