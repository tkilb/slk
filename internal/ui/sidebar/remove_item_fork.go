package sidebar

// RemoveItem drops the channel with the given ID from the sidebar (fork
// feature close-dm) and reports whether it was there. When the cursor was on
// the removed row it stays in that row's section: on the row that took its
// place, or, if it was the section's last row, on the row above it (the
// row above the section when nothing is left). Closing several DMs in a row neither
// jumps to the top nor into the next section.
func (m *Model) RemoveItem(id string) bool {
	idx := -1
	for i := range m.items {
		if m.items[i].ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return false
	}
	key, hadKey := m.currentCursorKey()
	onRemoved := hadKey && key.kind == navChannel && key.id == id
	oldCursor := m.cursor
	section := m.sectionFor(m.items[idx])

	items := make([]ChannelItem, 0, len(m.items)-1)
	items = append(items, m.items[:idx]...)
	items = append(items, m.items[idx+1:]...)
	m.items = items

	m.rebuildFilter()
	m.rebuildNavWithCursor(key, hadKey)
	if onRemoved {
		m.cursor = m.landingAfterRemoval(oldCursor, section)
	}
	m.cacheValid = false
	m.dirty()
	return true
}

// landingAfterRemoval picks the cursor row after a removal: the channel now at
// oldCursor if it is in the same section, otherwise the row above (the
// section's previous channel, or its header). When the section vanished
// with its last row, the row above where its header was.
func (m *Model) landingAfterRemoval(oldCursor int, section string) int {
	if oldCursor >= 0 && oldCursor < len(m.nav) {
		if n := m.nav[oldCursor]; n.kind == navChannel &&
			n.fi >= 0 && n.fi < len(m.filtered) &&
			m.sectionFor(m.items[m.filtered[n.fi]]) == section {
			return oldCursor
		}
	}
	for _, n := range m.nav {
		if n.kind == navHeader && n.header == section {
			return max(0, min(oldCursor-1, len(m.nav)-1))
		}
	}
	// The section held only that row and is gone, header and all: step over
	// the vanished header to the row above where the section was.
	return max(0, min(oldCursor-2, len(m.nav)-1))
}
