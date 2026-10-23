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
type threadSearch struct {
	query   string
	terms   []string // folded, highlighted by the renderer
	matches []int    // selection positions (parentSelected or a reply index), ascending
	idx     int      // index into matches of the current match
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
	s := threadSearch{query: query, terms: terms}
	if m.parent.TS != "" && m.messageMatches(m.parent, terms) {
		s.matches = append(s.matches, parentSelected)
	}
	for i, r := range m.replies {
		if m.messageMatches(r, terms) {
			s.matches = append(s.matches, i)
		}
	}
	for i, pos := range s.matches {
		if pos >= m.selected {
			s.idx = i
			break
		}
	}
	m.search = s
	if len(s.matches) > 0 {
		m.selectPosition(s.matches[s.idx])
	}
	m.InvalidateCache()
	return len(s.matches)
}

// StepSearch moves to the next (dir > 0) or previous match, wrapping, and
// returns the 1-based position and the match count. With no matches it
// returns 0, 0.
func (m *Model) StepSearch(dir int) (idx, count int) {
	n := len(m.search.matches)
	if n == 0 {
		return 0, 0
	}
	step := 1
	if dir < 0 {
		step = -1
	}
	m.search.idx = (m.search.idx + step + n) % n
	m.selectPosition(m.search.matches[m.search.idx])
	return m.search.idx + 1, n
}

// SearchStatus returns the 1-based current match and the match count.
func (m *Model) SearchStatus() (idx, count int) {
	if len(m.search.matches) == 0 {
		return 0, 0
	}
	return m.search.idx + 1, len(m.search.matches)
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
	folded := text.Fold(messages.FlattenMrkdwn(msg.Text, resolve, nil))
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
