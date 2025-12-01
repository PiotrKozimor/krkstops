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

type Trie[T any] struct {
	root        *node[T]
	transformer transform.Transformer
	compare     func(T, T) int
	equal       func(T, T) bool
}

type node[T any] struct {
	children map[rune]*node[T]
	results  []T
}

type Entry[T any] struct {
	Item T
	Word string
}

func New[T any](compare func(T, T) int, equal func(T, T) bool) Trie[T] {
	return Trie[T]{
		root: &node[T]{
			results:  []T{},
			children: map[rune]*node[T]{},
		},
		compare: compare,
		equal:   equal,
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

func (t *Trie[T]) Insert(entries ...Entry[T]) {
	for _, entry := range entries {
		if len(entry.Word) == 0 {
			continue
		}

		entry.Word = t.mustNormalize(entry.Word)
		t.insert(entry)
	}
}

func (t *Trie[T]) insert(entry Entry[T]) {
	currentNode := t.root
	for index, character := range entry.Word {
		child, ok := currentNode.children[character]
		if !ok {
			child = new(node[T])
			child.children = make(map[rune]*node[T])
			child.results = make([]T, 0)
			currentNode.children[character] = child
		}
		if index >= minSearchLen-1 {
			child.results = append(child.results, entry.Item)
		}
		currentNode = child
	}
}

// InsertWords will split "Foo i Bar" into two entries
// - "Foo i Bar"
// - "Bar"
func (t *Trie[T]) InsertWords(entries ...Entry[T]) {
	for _, e := range entries {
		w := t.mustNormalize(e.Word)
		f := strings.Fields(w)
		for i := range f {
			if utf8.RuneCount([]byte(f[i])) > 1 {
				w := strings.Join(f[i:], " ")
				t.insert(Entry[T]{
					Item: e.Item,
					Word: w,
				})
			}
		}
	}
}

func (t *Trie[T]) mustNormalize(word string) string {
	normal, _, err := transform.String(t.transformer, word)
	if err != nil {
		panic(err)
	}
	return normal
}
