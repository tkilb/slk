package thread

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gammons/slk/internal/text"
	"github.com/gammons/slk/internal/ui/messages"
)

// threadSearch is the client-side search state (fork feature thread-search).
// Matching is local: it covers the parent and the replies already loaded.
// Only the query is stored; matches are recomputed on demand from the live
// replies, so deletes, edits and new replies never leave stale indices.
type threadSearch struct {
	query string
	terms []string // folded, highlighted by the renderer
}

// matchPositions returns the selection positions (parentSelected or a reply
// index, ascending) of every message matching the active terms.
func (m *Model) matchPositions() []int {
	var out []int
	if len(m.search.terms) == 0 {
		return out
	}
	if m.parent.TS != "" && m.messageMatches(m.parent, m.search.terms) {
		out = append(out, parentSelected)
	}
	for i, r := range m.replies {
		if m.messageMatches(r, m.search.terms) {
			out = append(out, i)
		}
	}
	return out
}

// SetSearch starts or updates the search and returns the match count. A
// message matches when every whitespace-separated term, folded for case and
// accents, is a word prefix of its flattened text (the rule the renderer uses
// to highlight). The cursor jumps to the nearest match at or after it,
// wrapping to the first; with no match it stays put. A blank query clears the
// search.
func (m *Model) SetSearch(query string) int {
	terms := strings.Fields(text.Fold(query))
	if len(terms) == 0 {
		m.ClearSearch()
		return 0
	}
	m.search = threadSearch{query: query, terms: terms}
	matches := m.matchPositions()
	if len(matches) > 0 {
		target := matches[0]
		for _, pos := range matches {
			if pos >= m.selected {
				target = pos
				break
			}
		}
		m.selectPosition(target)
	}
	m.InvalidateCache()
	return len(matches)
}

// StepSearch moves to the next (dir > 0) or previous match relative to the
// cursor, wrapping, and returns the 1-based position and the match count.
// With no matches it returns 0, 0.
func (m *Model) StepSearch(dir int) (idx, count int) {
	matches := m.matchPositions()
	if len(matches) == 0 {
		return 0, 0
	}
	target := matches[0] // forward wrap
	if dir < 0 {
		target = matches[len(matches)-1] // backward wrap
	}
	for _, pos := range matches {
		if dir > 0 && pos > m.selected {
			target = pos
			break
		}
		if dir < 0 && pos < m.selected {
			target = pos
		}
	}
	m.selectPosition(target)
	return m.SearchStatus()
}

// SearchStatus returns the cursor's 1-based position among the matches and
// the match count. idx is 0 when the cursor is not on a match.
func (m *Model) SearchStatus() (idx, count int) {
	matches := m.matchPositions()
	for i, pos := range matches {
		if pos == m.selected {
			return i + 1, len(matches)
		}
	}
	return 0, len(matches)
}

// SearchQuery returns the active query as typed ("" when none).
func (m *Model) SearchQuery() string { return m.search.query }

// SearchActive reports whether a search, even one with no matches, is set.
func (m *Model) SearchActive() bool { return m.search.query != "" }

// ClearSearch drops the search and its highlights.
func (m *Model) ClearSearch() {
	if m.search.query == "" && len(m.search.terms) == 0 {
		return
	}
	m.search = threadSearch{}
	m.InvalidateCache()
}

func (m *Model) selectPosition(pos int) {
	if m.selected != pos {
		m.selected = pos
		m.InvalidateCache()
	}
}

func (m *Model) messageMatches(msg messages.MessageItem, terms []string) bool {
	resolve := func(id string) (string, bool) {
		name, ok := m.userNames[id]
		return name, ok
	}
	// The text the renderer draws (rich_text blocks reconstructed), with
	// user groups resolved, so what the user sees is what matches.
	folded := text.Fold(messages.FlattenMrkdwnWithUserGroups(messages.MessageTextSource(msg), resolve, nil, m.userGroups))
	for _, t := range terms {
		if !hasWordPrefix(folded, t) {
			return false
		}
	}
	return true
}

// hasWordPrefix reports whether term occurs in s at the start of a word:
// at the start of s or after a rune that is neither a letter nor a digit.
func hasWordPrefix(s, term string) bool {
	for from := 0; from <= len(s); {
		j := strings.Index(s[from:], term)
		if j < 0 {
			return false
		}
		j += from
		if j == 0 {
			return true
		}
		if r, _ := utf8.DecodeLastRuneInString(s[:j]); !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
		_, size := utf8.DecodeRuneInString(s[j:])
		from = j + size
	}
	return false
}
