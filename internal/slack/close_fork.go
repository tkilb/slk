package slackclient

import (
	"context"
	"encoding/json"
	"fmt"
)

// ConversationClosedHandler is optionally implemented by an EventHandler that
// wants to hear when a DM or group DM is closed, from this client or any
// other (im_close, mpim_close). Handlers without it ignore those events.
// Fork feature close-dm.
type ConversationClosedHandler interface {
	OnConversationClosed(channelID string)
}

// wsConversationClosedEvent is the shape of im_close and mpim_close.
type wsConversationClosedEvent struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
}

// dispatchConversationClosed forwards an im_close or mpim_close frame to the
// handler when it implements ConversationClosedHandler.
func dispatchConversationClosed(data []byte, handler EventHandler) {
	h, ok := handler.(ConversationClosedHandler)
	if !ok {
		return
	}
	var evt wsConversationClosedEvent
	if err := json.Unmarshal(data, &evt); err != nil || evt.Channel == "" {
		return
	}
	h.OnConversationClosed(evt.Channel)
}

// CloseConversation closes a DM or group DM (conversations.close), hiding it
// from the user's sidebar until a new message reopens it. alreadyClosed is
// true when Slack reports it was not open.
func (c *Client) CloseConversation(ctx context.Context, channelID string) (alreadyClosed bool, err error) {
	if channelID == "" {
		return false, fmt.Errorf("closing conversation: channel ID required")
	}
	_, alreadyClosed, err = c.api.CloseConversationContext(ctx, channelID)
	if err != nil {
		return false, fmt.Errorf("closing conversation: %w", err)
	}
	return alreadyClosed, nil
}
