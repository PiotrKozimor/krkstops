package trie

func (t *Trie[T]) SearchInDistance(term string, distance int) []T {
	term = t.mustNormalize(term)

	runes := []rune(term)
	if len(runes) < 2 {
		return nil
	}

	return searchNode(t.root, runes, distance, 0, true)
}

func (t *Trie[T]) SearchWithinDistance(term string, maxDistance int) []T {
	term = t.mustNormalize(term)

	runes := []rune(term)
	if len(runes) < 2 {
		return nil
	}

	return searchNode(t.root, runes, maxDistance, 0, false)
}

func searchNode[T any](n *node[T], r []rune, distance, level int, exact bool) []T {
	if len(r) == 0 {
		if !exact || distance == 0 {
			return n.results
		} else {
			return nil
		}
	}
	if distance == 0 {
		if node, ok := n.children[r[0]]; ok {
			return searchNode(node, r[1:], distance, level+1, exact)
		}
	} else {
		results := []T{}
		for nodeR, node := range n.children {
			if d := chebyshevDistance(nodeR, r[0]); d <= distance {
				r := searchNode(node, r[1:], distance-d, level+1, exact)
				results = append(results, r...)
			}
		}
		return results
	}
	return nil
}
