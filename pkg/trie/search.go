package trie

import "slices"

func (s *Trie[T]) Search(term string, limit int) []T {
	items := s.SearchExact(term)
	slices.SortFunc(items, s.compare)
	if len(items) >= limit {
		return items[:limit]
	}

	extra := s.SearchInDistance(term, 1)
	slices.SortFunc(extra, s.compare)
	items = append(items, extra...)
	if len(items) > limit {
		return items[:limit]
	}
	return items
}
