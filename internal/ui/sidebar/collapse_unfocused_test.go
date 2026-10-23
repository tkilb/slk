package sidebar

import (
	"strings"
	"testing"
)

func unfocusedCollapseModel() *Model {
	m := New([]ChannelItem{
		{ID: "C1", Name: "general", Type: "channel"},
		{ID: "C2", Name: "random", Type: "channel"},
		{ID: "D1", Name: "alice", Type: "dm"},
		{ID: "D2", Name: "bob", Type: "dm"},
	})
	m.SetCollapseWhenUnfocused(true)
	return &m
}

func TestCollapseWhenUnfocused_HidesRowsExceptActive(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetActiveChannelID("D2")
	m.SetFocused(false)

	view := m.View(30, 30)
	for _, hidden := range []string{"general", "random", "alice"} {
		if strings.Contains(view, hidden) {
			t.Errorf("unfocused sidebar should hide %q:\n%s", hidden, view)
		}
	}
	if !strings.Contains(view, "bob") {
		t.Errorf("unfocused sidebar should keep the active channel visible:\n%s", view)
	}
	if !strings.Contains(view, "▸") {
		t.Errorf("unfocused sections should show the collapsed glyph:\n%s", view)
	}
}

func TestCollapseWhenUnfocused_FocusExpandsEverySection(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(false)
	m.SetFocused(true)

	view := m.View(30, 30)
	// "Channels" is collapsed by default; focus overrides that.
	for _, shown := range []string{"general", "random", "alice", "bob"} {
		if !strings.Contains(view, shown) {
			t.Errorf("focused sidebar should show %q:\n%s", shown, view)
		}
	}
}

func TestCollapseWhenUnfocused_ToggleIsNoOp(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.MoveDown() // Activity
	m.MoveDown() // first section header
	if _, ok := m.IsSectionHeaderSelected(); !ok {
		t.Fatalf("expected a section header under the cursor")
	}
	before := m.View(30, 30)
	m.ToggleCollapseSelected()
	if after := m.View(30, 30); after != before {
		t.Errorf("toggling a section while the mode is on should change nothing")
	}
}

func TestCollapseWhenUnfocused_CursorReturnsToActiveRow(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetActiveChannelID("D2")
	m.SetFocused(true)
	m.SelectByID("C1") // browse to a row that will be hidden on blur
	if m.SelectedID() != "C1" {
		t.Fatalf("precondition: C1 selected, got %q", m.SelectedID())
	}

	m.SetFocused(false)
	m.SetFocused(true)
	if got := m.SelectedID(); got != "D2" {
		t.Errorf("cursor should return to the active channel, got %q", got)
	}
}

func TestCollapseWhenUnfocused_HeaderCursorSurvivesBlur(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.MoveDown() // Activity
	m.MoveDown() // header
	want, ok := m.IsSectionHeaderSelected()
	if !ok {
		t.Fatalf("expected a header under the cursor")
	}

	m.SetFocused(false)
	m.SetFocused(true)
	if got, ok := m.IsSectionHeaderSelected(); !ok || got != want {
		t.Errorf("header cursor lost across blur/focus: got %q ok=%v want %q", got, ok, want)
	}
}

func TestCollapseWhenUnfocused_DisabledLeavesBehaviorAlone(t *testing.T) {
	m := New([]ChannelItem{
		{ID: "D1", Name: "alice", Type: "dm"},
	})
	m.SetFocused(false)
	if !strings.Contains(m.View(20, 30), "alice") {
		t.Errorf("without the mode, an unfocused sidebar must not collapse")
	}
}

func TestCollapseWhenUnfocused_CursorFollowsNewActiveChannel(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("D1")
	m.SetActiveChannelID("D1")
	m.SetFocused(false)        // cursor stays on the visible, active D1 row
	m.SetActiveChannelID("D2") // switched from the finder while unfocused
	m.SetFocused(true)
	if got := m.SelectedID(); got != "D2" {
		t.Errorf("cursor should land on the new active channel, got %q", got)
	}
}

func TestCollapseWhenUnfocused_ToggleSelectedReportsNoOp(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.MoveDown()
	m.MoveDown()
	if m.ToggleCollapseSelected() {
		t.Errorf("ToggleCollapseSelected must report false when it does nothing")
	}
}
