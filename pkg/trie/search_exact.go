package trie

import "unicode/utf8"

// Search will return nil if not results were found
func (t *Trie) SearchExact(term string) []uint {
	term = t.mustNormalize(term)
	currentNode := t.root

	if utf8.RuneCountInString(term) < 2 {
		return nil
	}

	for _, character := range term {
		if node, ok := currentNode.children[character]; ok {
			currentNode = node
		} else {
			return nil
		}
	}

	return currentNode.results
}
