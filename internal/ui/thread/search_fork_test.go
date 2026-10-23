package thread

import (
	"testing"

	"github.com/gammons/slk/internal/ui/messages"
)

func searchThread() *Model {
	m := New()
	parent := messages.MessageItem{TS: "10.0", UserID: "U1", UserName: "alice", Text: "Deploy plan for Friday", ThreadTS: "10.0"}
	replies := []messages.MessageItem{
		{TS: "11.0", UserID: "U2", UserName: "bob", Text: "looks fine"},
		{TS: "12.0", UserID: "U1", UserName: "alice", Text: "deployment window is <https://x.example|tight>"},
		{TS: "13.0", UserID: "U2", UserName: "bob", Text: "Café menu unrelated"},
		{TS: "14.0", UserID: "U1", UserName: "alice", Text: "ship the DEPLOY"},
	}
	m.SetThread(parent, replies, "C1", "10.0")
	return m
}

func selectedTS(m *Model) string {
	r := m.SelectedReply()
	if r == nil {
		return ""
	}
	return r.TS
}

func TestSetSearch_CountsWordPrefixMatchesIncludingParent(t *testing.T) {
	m := searchThread()
	if got := m.SetSearch("deploy"); got != 3 { // parent, reply 12, reply 14
		t.Fatalf("matches = %d, want 3", got)
	}
}

func TestSetSearch_FoldsCaseAndAccents(t *testing.T) {
	m := searchThread()
	if got := m.SetSearch("cafe"); got != 1 {
		t.Fatalf("matches = %d, want 1", got)
	}
	if selectedTS(m) != "13.0" {
		t.Errorf("selection = %q, want 13.0", selectedTS(m))
	}
}

func TestSetSearch_MatchesLinkLabelNotURL(t *testing.T) {
	m := searchThread()
	if got := m.SetSearch("tight"); got != 1 {
		t.Errorf("link label should match, got %d", got)
	}
	if got := m.SetSearch("example"); got != 0 {
		t.Errorf("link URL should not match, got %d", got)
	}
}

func TestSetSearch_AllTermsMustMatch(t *testing.T) {
	m := searchThread()
	if got := m.SetSearch("deploy friday"); got != 1 {
		t.Errorf("matches = %d, want only the parent", got)
	}
}

func TestSetSearch_MidWordIsNotAMatch(t *testing.T) {
	m := searchThread()
	if got := m.SetSearch("eploy"); got != 0 {
		t.Errorf("mid-word text should not match, got %d", got)
	}
}

func TestSetSearch_SelectsNearestMatchAtOrAfterCursor(t *testing.T) {
	m := searchThread()
	m.SelectByIndex(1) // reply 12.0, itself a match
	m.SetSearch("deploy")
	if selectedTS(m) != "12.0" {
		t.Errorf("selection = %q, want 12.0", selectedTS(m))
	}
	m.SelectByIndex(2) // 13.0, not a match; next is 14.0
	m.SetSearch("deplo")
	if selectedTS(m) != "14.0" {
		t.Errorf("selection = %q, want 14.0", selectedTS(m))
	}
}

func TestStepSearch_WrapsBothWays(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy") // cursor starts at the newest reply (14.0)
	if i, n := m.StepSearch(1); i != 1 || n != 3 || selectedTS(m) != "10.0" {
		t.Errorf("forward wrap: idx=%d n=%d ts=%q", i, n, selectedTS(m))
	}
	if i, _ := m.StepSearch(-1); i != 3 || selectedTS(m) != "14.0" {
		t.Errorf("backward wrap: idx=%d ts=%q", i, selectedTS(m))
	}
}

func TestSetSearch_NoMatchLeavesSelection(t *testing.T) {
	m := searchThread()
	before := selectedTS(m)
	if got := m.SetSearch("zzz"); got != 0 {
		t.Fatalf("matches = %d, want 0", got)
	}
	if selectedTS(m) != before {
		t.Errorf("selection moved to %q", selectedTS(m))
	}
	if i, n := m.StepSearch(1); i != 0 || n != 0 {
		t.Errorf("step with no matches: idx=%d n=%d", i, n)
	}
}

func TestClearSearchAndReset(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy")
	m.ClearSearch()
	if m.SearchActive() || m.SearchQuery() != "" {
		t.Errorf("search should be inactive after ClearSearch")
	}
	m.SetSearch("deploy")
	m.SetThread(messages.MessageItem{TS: "1.0"}, nil, "C2", "1.0")
	if m.SearchActive() {
		t.Errorf("opening another thread should drop the search")
	}
}

func TestSetSearch_EmptyQueryIsInactive(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy")
	m.SetSearch("  ")
	if m.SearchActive() {
		t.Errorf("a blank query should clear the search")
	}
}

func TestStepSearch_SurvivesReplyRemoval(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy") // parent, 12.0, 14.0
	m.RemoveMessageByTS("12.0")
	if i, n := m.StepSearch(-1); n != 2 || selectedTS(m) != "10.0" || i != 1 {
		t.Errorf("after removal: idx=%d n=%d ts=%q, want 1/2 on the parent", i, n, selectedTS(m))
	}
}

func TestSetThread_SameThreadKeepsSearch(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy")
	m.SetThread(m.ParentMsg(), m.Replies(), "C1", "10.0") // a refresh
	if !m.SearchActive() || m.SearchQuery() != "deploy" {
		t.Errorf("refreshing the open thread must keep its search")
	}
}

func TestSearchStatus_ZeroWhenCursorOffMatches(t *testing.T) {
	m := searchThread()
	m.SetSearch("deploy")
	m.SelectByIndex(0) // 11.0, not a match
	if i, n := m.SearchStatus(); i != 0 || n != 3 {
		t.Errorf("status = %d/%d, want 0/3", i, n)
	}
}
