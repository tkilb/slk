package slackclient

import (
	"context"
	"errors"
	"testing"
)

type closedRecorder struct {
	mockEventHandler
	closed []string
}

func (r *closedRecorder) OnConversationClosed(channelID string) {
	r.closed = append(r.closed, channelID)
}

func TestCloseConversation_CallsSlackWithTheChannel(t *testing.T) {
	var got string
	c := &Client{api: &mockSlackAPI{closeConversationContextFn: func(_ context.Context, id string) (bool, bool, error) {
		got = id
		return false, false, nil
	}}}
	already, err := c.CloseConversation(context.Background(), "D1")
	if err != nil || already {
		t.Fatalf("already=%v err=%v, want false, nil", already, err)
	}
	if got != "D1" {
		t.Errorf("closed %q, want D1", got)
	}
}

func TestCloseConversation_ReportsAlreadyClosed(t *testing.T) {
	c := &Client{api: &mockSlackAPI{closeConversationContextFn: func(context.Context, string) (bool, bool, error) {
		return false, true, nil
	}}}
	if already, err := c.CloseConversation(context.Background(), "D1"); err != nil || !already {
		t.Errorf("already=%v err=%v, want true, nil", already, err)
	}
}

func TestCloseConversation_WrapsErrors(t *testing.T) {
	boom := errors.New("boom")
	c := &Client{api: &mockSlackAPI{closeConversationContextFn: func(context.Context, string) (bool, bool, error) {
		return false, false, boom
	}}}
	if _, err := c.CloseConversation(context.Background(), "D1"); !errors.Is(err, boom) {
		t.Errorf("err = %v, want it to wrap boom", err)
	}
}

func TestCloseConversation_RejectsEmptyID(t *testing.T) {
	c := &Client{api: &mockSlackAPI{}}
	if _, err := c.CloseConversation(context.Background(), ""); err == nil {
		t.Errorf("an empty channel ID must be rejected")
	}
}

func TestDispatch_IMAndMPIMCloseReachAnInterestedHandler(t *testing.T) {
	for _, payload := range []string{
		`{"type":"im_close","user":"U1","channel":"D1"}`,
		`{"type":"mpim_close","user":"U1","channel":"D1"}`,
	} {
		h := &closedRecorder{}
		dispatchWebSocketEvent([]byte(payload), h)
		if len(h.closed) != 1 || h.closed[0] != "D1" {
			t.Errorf("%s: closed = %v, want [D1]", payload, h.closed)
		}
	}
}

func TestDispatch_CloseIsIgnoredByHandlersWithoutTheInterface(t *testing.T) {
	h := &mockEventHandler{}
	dispatchWebSocketEvent([]byte(`{"type":"im_close","user":"U1","channel":"D1"}`), h) // must not panic
}

func TestDispatch_CloseWithoutAChannelIsIgnored(t *testing.T) {
	h := &closedRecorder{}
	dispatchWebSocketEvent([]byte(`{"type":"im_close","user":"U1"}`), h)
	if len(h.closed) != 0 {
		t.Errorf("closed = %v, want none", h.closed)
	}
}
