package trie

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	search := New(func(s1, s2 string) int {
		return strings.Compare(s1, s2)
	}, func(s1, s2 string) bool {
		return s1 == s2
	})
	search.InsertWords([]Entry[string]{
		{
			Item: "1",
			Word: "Starowiślna",
		},
	}...)

	assert.ElementsMatch(t, search.root.children['s'].children['t'].children['a'].results, []string{"1"})
	results := search.Search("ata", 10)
	assert.Equal(t, "1", results[0])
}
