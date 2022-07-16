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
		char     string
		isWord   bool
		Children [SIZE]*Node
	}

	Dictionary struct {
		RootNode *Node
	}
)

func newNode(char string) *Node {
	return &Node{
		char:   char,
		isWord: false,
	}
}

func newDictionary() *Dictionary {
	trie := &Dictionary{
		RootNode: newNode(ROOT),
	}

	for word := range dictionary {
		trie.insert(word)
	}

	return trie
}

func (d *Dictionary) insert(word string) {
	current := d.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	for i := 0; i < len(strippedWord); i++ {
		char := string(strippedWord[i])
		index := strippedWord[i] - 'a'

		if current.Children[index] == nil {
			current.Children[index] = newNode(char)
		}

		current = current.Children[index]
	}
	current.isWord = true
}

func (d *Dictionary) searchWord(word string) bool {
	current := d.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	for i := 0; i < len(strippedWord); i++ {
		index := strippedWord[i] - 'a'

		if current == nil || current.Children[index] == nil {
			return false
		}
		current = current.Children[index]
	}

	return current.isWord
}

func (n *Node) IsWord() bool {
	return n.isWord
}
