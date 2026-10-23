package sidebar

// RemoveItem drops the channel with the given ID from the sidebar (fork
// feature close-dm) and reports whether it was there. When the cursor was on
// the removed row it moves to the row that took its place, so closing several
// DMs in a row does not send the cursor back to the top.
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

	items := make([]ChannelItem, 0, len(m.items)-1)
	items = append(items, m.items[:idx]...)
	items = append(items, m.items[idx+1:]...)
	m.items = items

	m.rebuildFilter()
	m.rebuildNavWithCursor(key, hadKey)
	if onRemoved {
		m.cursor = min(oldCursor, len(m.nav)-1)
	}
	m.cacheValid = false
	m.dirty()
	return true
}
