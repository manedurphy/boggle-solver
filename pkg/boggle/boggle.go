package boggle

import (
	"errors"
	"fmt"
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
	"PLAY":        {},
}

type Boggle struct {
	found   map[string]struct{}
	board   [][]string
	visited [][]bool
	m       int
	n       int
}

func New(board [][]string) (*Boggle, error) {
	if len(board) == 0 {
		return nil, errors.New("board is invalid")
	}

	m := len(board)
	n := len(board[0])

	visited := make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	return &Boggle{
		found:   map[string]struct{}{},
		board:   board,
		visited: visited,
		m:       m,
		n:       n,
	}, nil
}

func (b *Boggle) Solve() []string {
	fmt.Println("solving boggle board...")
	runningString := ""

	for i := 0; i < b.m; i++ {
		for j := 0; j < b.n; j++ {
			b.findWords(i, j, runningString)
		}
	}

	fmt.Println("board solved!")
	return b.getFoundList()
}

func (b *Boggle) findWords(i, j int, runningString string) {
	b.visited[i][j] = true
	runningString += b.board[i][j]

	if _, valid := dictionary[runningString]; valid && len(runningString) >= 3 {
		b.found[runningString] = struct{}{}
	}

	for row := i - 1; row <= i+1 && row < len(b.board); row++ {
		for col := j - 1; col <= j+1 && col < len(b.board[0]); col++ {
			if row >= 0 && col >= 0 && !b.visited[row][col] {
				b.findWords(row, col, runningString)
			}
		}
	}

	b.visited[i][j] = false
}

func (b *Boggle) getFoundList() []string {
	var foundList []string

	for found := range b.found {
		foundList = append(foundList, found)
	}

	return foundList
}
