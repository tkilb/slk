package main

import (
	"testing"

	"github.com/gammons/slk/internal/core"
	slackclient "github.com/gammons/slk/internal/slack"
	"github.com/gammons/slk/internal/ui/sidebar"
)

var _ slackclient.ConversationClosedHandler = (*rtmEventHandler)(nil)

func closeTestHandler() (*rtmEventHandler, *WorkspaceContext) {
	wctx := &WorkspaceContext{
		UserNames:         newUserNameStore(nil),
		UserNamesByHandle: map[string]string{},
		Channels: []sidebar.ChannelItem{
			{ID: "C1", Name: "general", Type: "channel"},
			{ID: "D1", Name: "alice", Type: "dm"},
		},
		FinderItems: []core.ChannelFinderItem{
			{ID: "C1", Name: "general", Type: "channel", Joined: true},
			{ID: "D1", Name: "alice", Type: "dm"},
		},
	}
	return &rtmEventHandler{wsCtx: wctx, workspaceID: "T1"}, wctx
}

func TestOnConversationClosed_DropsTheSidebarItemButKeepsTheFinderEntry(t *testing.T) {
	h, wctx := closeTestHandler()
	h.OnConversationClosed("D1")
	if len(wctx.Channels) != 1 || wctx.Channels[0].ID != "C1" {
		t.Errorf("Channels = %+v, want only C1", wctx.Channels)
	}
	// The finder offers "message this person" for closed DMs, so it keeps them.
	if len(wctx.FinderItems) != 2 {
		t.Errorf("FinderItems = %d, want 2 (closed DMs stay findable)", len(wctx.FinderItems))
	}
}

func TestOnConversationClosed_NeverClosesRegularChannels(t *testing.T) {
	h, wctx := closeTestHandler()
	h.OnConversationClosed("C1")
	if len(wctx.Channels) != 2 {
		t.Errorf("a regular channel must stay, Channels = %+v", wctx.Channels)
	}
}

func TestOnConversationClosed_UnknownIDAndNilContextAreSafe(t *testing.T) {
	h, wctx := closeTestHandler()
	h.OnConversationClosed("nope")
	if len(wctx.Channels) != 2 {
		t.Errorf("unknown ID changed Channels: %+v", wctx.Channels)
	}
	(&rtmEventHandler{}).OnConversationClosed("D1") // wsCtx nil: must not panic
}
