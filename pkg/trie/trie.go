package trie

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const minSearchLen = 2

type Trie struct {
	root        *node
	transformer transform.Transformer
}

type node struct {
	children map[rune]*node
	results  []uint
}

type Entry struct {
	Id   uint
	Word string
}

func New() Trie {
	return Trie{
		root: &node{
			results:  []uint{},
			children: map[rune]*node{},
		},
		transformer: transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
			runes.Remove(runes.In(unicode.Punct)),
			runes.Map(unicode.ToLower),
			runes.Map(func(r rune) rune {
				if r == 'ł' {
					return 'l'
				}
				return r
			}),
		),
	}
}

func (t *Trie) Insert(entries ...Entry) {
	for _, entry := range entries {
		if len(entry.Word) == 0 {
			continue
		}

		entry.Word = t.mustNormalize(entry.Word)
		t.insert(entry)
	}
}

func (t *Trie) insert(entry Entry) {
	currentNode := t.root
	for index, character := range entry.Word {
		child, ok := currentNode.children[character]
		if !ok {
			child = new(node)
			child.children = make(map[rune]*node)
			child.results = make([]uint, 0)
			currentNode.children[character] = child
		}
		if index >= minSearchLen-1 {
			child.results = append(child.results, entry.Id)
		}
		currentNode = child
	}
}

// InsertWords will split "Foo i Bar" into two entries
// - "Foo i Bar"
// - "Bar"
func (t *Trie) InsertWords(entries ...Entry) {
	for _, e := range entries {
		w := t.mustNormalize(e.Word)
		f := strings.Fields(w)
		for i := range f {
			if utf8.RuneCount([]byte(f[i])) > 1 {
				w := strings.Join(f[i:], " ")
				t.insert(Entry{
					Id:   e.Id,
					Word: w,
				})
			}
		}
	}
}

func (t *Trie) mustNormalize(word string) string {
	normal, _, err := transform.String(t.transformer, word)
	if err != nil {
		panic(err)
	}
	return normal
}
