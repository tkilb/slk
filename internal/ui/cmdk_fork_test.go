package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// Cmd+K arrives as super+k in terminals that report the super modifier
// (kitty keyboard protocol) and opens the channel finder.
func TestCmdKOpensChannelFinder(t *testing.T) {
	runKeyCases(t, ModeNormal, []keyCase{{
		name:     "super+k opens the channel finder",
		opts:     normalOpts(),
		key:      keyMod('k', tea.ModSuper),
		wantMode: ModeChannelFinder,
		assert: func(t *testing.T, a *App, _ tea.Cmd) {
			if !a.channelFinder.IsVisible() {
				t.Error("channel finder not visible")
			}
		},
	}})
}
