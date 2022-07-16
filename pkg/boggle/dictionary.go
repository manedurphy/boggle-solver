package boggle

import (
	"strings"
)

const (
	SIZE = 26
	ROOT = "root"
)

var dictionary = map[string]struct{}{
	"DATA":        {},
	"SCIENCE":     {},
	"SOFTWARE":    {},
	"ENGINEERING": {},
	"GOLANG":      {},
	"RUST":        {},
	"PLEASE":      {},
	"EXCUSE":      {},
	"MY":          {},
	"DEAR":        {},
	"AUNT":        {},
	"SALLY":       {},
	"BOGGLE":      {},
	"PLEE":        {},
	"PLAY":        {},
}

type (
	Node struct {
		Char     string
		IsWord   bool
		Children [SIZE]*Node
	}

	Trie struct {
		RootNode *Node
	}
)

func NewNode(char string) *Node {
	return &Node{
		Char:   char,
		IsWord: false,
	}
}

func NewTrie() *Trie {
	trie := &Trie{
		RootNode: NewNode(ROOT),
	}

	for word := range dictionary {
		trie.Insert(word)
	}

	return trie
}

func (t *Trie) Insert(word string) {
	current := t.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	for i := 0; i < len(strippedWord); i++ {
		char := string(strippedWord[i])
		index := strippedWord[i] - 'a'

		if current.Children[index] == nil {
			current.Children[index] = NewNode(char)
		}

		current = current.Children[index]
	}
	current.IsWord = true
}

func (t *Trie) SearchWord(word string) bool {
	current := t.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	for i := 0; i < len(strippedWord); i++ {
		index := strippedWord[i] - 'a'

		if current == nil || current.Children[index] == nil {
			return false
		}
		current = current.Children[index]
	}

	return current.IsWord
}

/*
func (t *Trie) SearchWord(i, j int, word string, visited [][]bool) bool {
	current := t.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	if t.isSafe(i, j, 3, 3, visited) {
		visited[i][j] = true

		for k := 0; k < SIZE; k++ {
			if t.RootNode.Children[k] != nil {
				if t.isSafe(i+1, j+1, 3, 3, visited) {
					strippedWord +
				}
			}
		}

		for i := 0; i < len(strippedWord); i++ {
			index := strippedWord[i] - 'a'

			if current == nil || current.Children[index] == nil {
				return false
			}
			current = current.Children[index]
		}
	}

	return true
}

func (t *Trie) isSafe(i, j, m, n int, visited [][]bool) bool {
	return i >= 0 && i < m && j >= 0 && j < n && !visited[i][j]
}
*/
