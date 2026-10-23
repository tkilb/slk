package sidebar

import (
	"strings"
	"testing"
)

// Nav order in unfocusedCollapseModel: Threads, Activity, Direct Messages
// (D1, D2), Channels (C1, C2).

func TestJumpSection_NextFromChildLandsOnNextHeader(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("D1")
	if !m.JumpSection(1) {
		t.Fatalf("expected a next section")
	}
	if name, ok := m.IsSectionHeaderSelected(); !ok || name != "Channels" {
		t.Errorf("selected header = %q ok=%v, want Channels", name, ok)
	}
}

func TestJumpSection_NextFromThreadsLandsOnFirstHeader(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	if !m.JumpSection(1) {
		t.Fatalf("expected a first section")
	}
	if name, ok := m.IsSectionHeaderSelected(); !ok || name != "Direct Messages" {
		t.Errorf("selected header = %q ok=%v, want Direct Messages", name, ok)
	}
}

func TestJumpSection_PreviousFromChildLandsOnPreviousSectionHeader(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("C2")
	if !m.JumpSection(-1) {
		t.Fatalf("expected a previous section")
	}
	if name, ok := m.IsSectionHeaderSelected(); !ok || name != "Direct Messages" {
		t.Errorf("selected header = %q ok=%v, want Direct Messages", name, ok)
	}
}

func TestJumpSection_PreviousFromHeaderSkipsToTheOneBefore(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("C1")
	m.JumpSection(-1) // Direct Messages header
	m.JumpSection(1)  // Channels header
	if !m.JumpSection(-1) {
		t.Fatalf("expected a previous section")
	}
	if name, _ := m.IsSectionHeaderSelected(); name != "Direct Messages" {
		t.Errorf("selected header = %q, want Direct Messages", name)
	}
}

func TestJumpSection_StopsAtTheEnds(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	if m.JumpSection(-1) {
		t.Errorf("no previous section from Threads")
	}
	m.SelectByID("C2")
	if m.JumpSection(1) {
		t.Errorf("no next section after the last one")
	}
	if m.SelectedID() != "C2" {
		t.Errorf("a failed jump must not move the cursor, got %q", m.SelectedID())
	}
}

func TestJumpSection_OpensTheTargetSection(t *testing.T) {
	m := unfocusedCollapseModel()
	m.SetFocused(true)
	m.SelectByID("D1")
	m.JumpSection(1)
	view := m.View(30, 30)
	if !strings.Contains(view, "general") || strings.Contains(view, "alice") {
		t.Errorf("jumping should open Channels and close Direct Messages:\n%s", view)
	}
}
