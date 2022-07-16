package boggle

import (
	"errors"
	"fmt"
	"strings"
)

type Boggle struct {
	board      [][]string
	visited    [][]bool
	found      map[string]struct{}
	m          int
	n          int
	dictionary *Trie
}

func New(board [][]string) (*Boggle, error) {
	if len(board) == 0 {
		return nil, errors.New("board is invalid")
	}

	m := len(board)
	n := len(board[0])

	for idx, row := range board {
		for jdx, val := range row {
			board[idx][jdx] = strings.ToLower(val)
		}
	}

	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	return &Boggle{
		found:      make(map[string]struct{}),
		board:      board,
		visited:    visited,
		dictionary: NewTrie(),
		m:          m,
		n:          n,
	}, nil
}

func (b *Boggle) Solve() []string {
	fmt.Println("solving boggle board...")
	runningString := ""

	for i := 0; i < b.m; i++ {
		for j := 0; j < b.n; j++ {
			str := b.board[i][j]
			index := []byte(str)[0] - 'a'
			if node := b.dictionary.RootNode.Children[index]; node != nil {
				b.findWords(i, j, runningString+str, node)
			}
		}
	}

	fmt.Println("board solved!")
	return b.getFoundList()
}

func (b *Boggle) findWords(i, j int, runningString string, node *Node) {
	if node.IsWord && len(runningString) >= 3 {
		b.found[runningString] = struct{}{}
	}

	if b.isSafe(i, j) {
		b.visited[i][j] = true

		for k := 0; k < SIZE; k++ {
			if node.Children[k] != nil {
				ch := string(byte(k + 'a'))

				if b.isSafe(i+1, j+1) && b.board[i+1][j+1] == ch {
					b.findWords(i+1, j+1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i, j+1) && b.board[i][j+1] == ch {
					b.findWords(i, j+1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i-1, j+1) && b.board[i-1][j+1] == ch {
					b.findWords(i-1, j+1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i+1, j) && b.board[i+1][j] == ch {
					b.findWords(i+1, j, runningString+ch, node.Children[k])
				}
				if b.isSafe(i+1, j-1) && b.board[i+1][j-1] == ch {
					b.findWords(i+1, j-1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i, j-1) && b.board[i][j-1] == ch {
					b.findWords(i, j-1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i-1, j-1) && b.board[i-1][j-1] == ch {
					b.findWords(i-1, j-1, runningString+ch, node.Children[k])
				}
				if b.isSafe(i-1, j) && b.board[i-1][j] == ch {
					b.findWords(i-1, j, runningString+ch, node.Children[k])
				}
			}
		}
		b.visited[i][j] = false
	}
}

func (b *Boggle) getFoundList() []string {
	var foundList []string

	for found := range b.found {
		foundList = append(foundList, found)
	}

	return foundList
}

func (b *Boggle) isSafe(i, j int) bool {
	return i >= 0 && i < b.m && j >= 0 && j < b.n && !b.visited[i][j]
}
