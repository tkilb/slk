// internal/ui/close_conversation.go
//
// Fork feature close-dm: `x` in the sidebar and `/close` close a DM or group
// DM (Slack's conversations.close). The row leaves the sidebar once Slack
// confirms; a later message reopens the conversation. Closes made from any
// Slack client arrive as ConversationClosedMsg (im_close / mpim_close).
package ui

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"

	"github.com/gammons/slk/internal/ids"
)

var forkCloseConversation = key.NewBinding(key.WithKeys("x"), key.WithHelp("x", "close DM (sidebar)"))

const closeToastDuration = 3 * time.Second

// ConversationClosedResultMsg is the result of ChannelService.CloseConversation.
type ConversationClosedResultMsg struct {
	ChannelID string
	Err       error
}

// ConversationClosedMsg is sent when Slack reports a DM or group DM closed
// (im_close, mpim_close). Only the active workspace's sidebar is touched.
type ConversationClosedMsg struct {
	TeamID    string
	ChannelID string
}

// isClosableType reports whether a sidebar item type can be closed.
func isClosableType(t string) bool { return t == "dm" || t == "group_dm" }

// closeSelectedConversation closes the DM under the sidebar cursor. Rows that
// are not DMs get a toast; headers and the Threads and Activity rows do nothing.
func (a *App) closeSelectedConversation() tea.Cmd {
	item, ok := a.sidebar.SelectedItem()
	if !ok {
		return nil
	}
	if !isClosableType(item.Type) {
		return toastWithClear(a, "Only DMs and group DMs can be closed", closeToastDuration)
	}
	return teaCmd(a.channels.CloseConversation(ids.ChannelID(item.ID)))
}

// conversationName is the sidebar name for a channel ID, or "" if unknown.
func (a *App) conversationName(channelID string) string {
	for _, it := range a.sidebar.Items() {
		if it.ID == channelID {
			return it.Name
		}
	}
	return ""
}

// slashClose closes the conversation the command was typed in.
func slashClose(a *App, args commandArgs) tea.Cmd {
	chID := args.Origin.channelID
	for _, it := range a.sidebar.Items() {
		if it.ID == chID && isClosableType(it.Type) {
			return teaCmd(a.channels.CloseConversation(ids.ChannelID(chID)))
		}
	}
	return a.refuseCommand(args.Origin, "/close only works in a DM or group DM")
}

var reduceCloseConversation reducerFunc = func(a *App, msg tea.Msg) (tea.Cmd, bool) {
	switch m := msg.(type) {
	case ConversationClosedResultMsg:
		name := a.conversationName(m.ChannelID)
		if name == "" {
			name = "conversation"
		}
		if m.Err != nil {
			return toastWithClear(a, "Couldn't close "+name+": "+m.Err.Error(), closeToastDuration), true
		}
		a.sidebar.RemoveItem(m.ChannelID)
		return toastWithClear(a, "Closed "+name, closeToastDuration), true

	case ConversationClosedMsg:
		if m.TeamID == a.activeTeamID {
			a.sidebar.RemoveItem(m.ChannelID)
		}
		return nil, true
	}
	return nil, false
}
