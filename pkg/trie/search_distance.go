package trie

func (t *Trie) SearchInDistance(term string, distance int) []uint {
	term = t.mustNormalize(term)

	runes := []rune(term)
	if len(runes) < 2 {
		return nil
	}

	return searchNode(t.root, runes, distance, true)
}

func (t *Trie) SearchWithinDistance(term string, maxDistance int) []uint {
	term = t.mustNormalize(term)

	runes := []rune(term)
	if len(runes) < 2 {
		return nil
	}

	return searchNode(t.root, runes, maxDistance, false)
}

func searchNode(n *node, r []rune, distance int, exact bool) []uint {
	if len(r) == 0 {
		if !exact || distance == 0 {
			return n.results
		} else {
			return nil
		}
	}
	if distance == 0 {
		if node, ok := n.children[r[0]]; ok {
			return searchNode(node, r[:1], distance, exact)
		}
	} else {
		results := []uint{}
		for nodeR, node := range n.children {
			if d := chebyshevDistance(nodeR, r[0]); d <= distance {
				r := searchNode(node, r[1:], distance-d, exact)
				results = append(results, r...)
			}
		}
		return results
	}
	return nil
}
