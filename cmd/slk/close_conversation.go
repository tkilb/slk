package main

import (
	"context"
	"fmt"

	"github.com/gammons/slk/internal/core"
	"github.com/gammons/slk/internal/ids"
	"github.com/gammons/slk/internal/ui"
	"github.com/gammons/slk/internal/ui/sidebar"
)

// closeConversationCmd builds the ChannelService.CloseConversation closure
// (fork feature close-dm): it closes the DM or group DM on the active
// workspace's client and reports the outcome as ui.ConversationClosedResultMsg.
// The workspace context itself is updated by OnConversationClosed, when
// Slack's im_close / mpim_close event arrives on the serialized event owner.
func closeConversationCmd(ctx context.Context, active func() *WorkspaceContext) func(ids.ChannelID) core.Cmd {
	return func(channelID ids.ChannelID) core.Cmd {
		wctx := active()
		if wctx == nil || wctx.Client == nil {
			return func() core.Msg {
				return ui.ConversationClosedResultMsg{
					ChannelID: string(channelID),
					Err:       fmt.Errorf("no active workspace"),
				}
			}
		}
		client := wctx.Client
		return func() core.Msg {
			_, err := client.CloseConversation(ctx, string(channelID))
			return ui.ConversationClosedResultMsg{ChannelID: string(channelID), Err: err}
		}
	}
}

// OnConversationClosed handles im_close and mpim_close, whichever Slack client
// closed the conversation. It drops the DM or group DM from the workspace's
// sidebar list so a later workspace switch does not bring it back, and tells
// the UI when this workspace is on screen. The finder entry stays: the finder
// offers closed DMs as "message this person". Only DMs are ever removed.
func (h *rtmEventHandler) OnConversationClosed(channelID string) {
	if h.wsCtx == nil {
		return
	}
	kept := make([]sidebar.ChannelItem, 0, len(h.wsCtx.Channels))
	removed := false
	for _, it := range h.wsCtx.Channels {
		if it.ID == channelID && (it.Type == "dm" || it.Type == "group_dm") {
			removed = true
			continue
		}
		kept = append(kept, it)
	}
	if !removed {
		return
	}
	h.wsCtx.Channels = kept
	if h.program == nil || (h.isActive != nil && !h.isActive()) {
		return
	}
	h.program.Send(ui.ConversationClosedMsg{TeamID: h.workspaceID, ChannelID: channelID})
}
