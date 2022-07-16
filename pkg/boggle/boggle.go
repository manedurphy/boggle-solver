package boggle

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type (
	Boggle struct {
		dictionary *Dictionary
		board      [][]string
		visited    [][]bool
		score      int
		found      map[string]struct{}
		m          int
		n          int
	}

	Result struct {
		WordsFound []string `json:"words_found"`
		Score      int      `json:"score"`
	}
)

func New(board [][]string) (*Boggle, error) {
	var (
		m          int
		n          int
		visited    [][]bool
		dictionary *Dictionary
		err        error
	)

	if len(board) == 0 {
		return nil, errors.New("board is invalid")
	}

	m = len(board)
	n = len(board[0])

	if !isValidBoard(board, n) {
		return nil, errors.New("board is imbalanced")
	}

	for idx, row := range board {
		for jdx, val := range row {
			board[idx][jdx] = strings.ToLower(val)
		}
	}

	visited = make([][]bool, m)
	for i := range visited {
		visited[i] = make([]bool, n)
	}

	dictionary, err = newDictionary()
	if err != nil {
		return nil, fmt.Errorf("faild to create new dictionary: %w", err)
	}

	return &Boggle{
		found:      make(map[string]struct{}),
		board:      board,
		visited:    visited,
		dictionary: dictionary,
		m:          m,
		n:          n,
	}, nil
}

func (b *Boggle) Solve() *Result {
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

	return &Result{
		WordsFound: b.getFoundList(),
		Score:      b.score,
	}
}

func (b *Boggle) findWords(i, j int, runningString string, node *Node) {
	runningStringLength := len(runningString)
	if node.isWord && runningStringLength >= 3 {
		_, valid := b.found[runningString]

		if !valid {
			switch {
			case runningStringLength >= 3 && runningStringLength < 5:
				b.score += 1
			case runningStringLength == 5:
				b.score += 2
			case runningStringLength == 6:
				b.score += 3
			case runningStringLength == 7:
				b.score += 5
			default:
				b.score += 11
			}
		}
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
	sort.Strings(foundList)

	return foundList
}

func (b *Boggle) isSafe(i, j int) bool {
	return i >= 0 && i < b.m && j >= 0 && j < b.n && !b.visited[i][j]
}

func isValidBoard(board [][]string, n int) bool {
	for _, row := range board {
		if len(row) != n {
			return false
		}
	}

	return true
}
