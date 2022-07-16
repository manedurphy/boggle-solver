package boggle

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"
)

const (
	SIZE = 26
	ROOT = "root"
)

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

func newDictionary(wordsFilePath string) (*Dictionary, error) {
	var (
		trie      *Dictionary
		zipReader *zip.ReadCloser
		err       error
	)

	trie = &Dictionary{
		RootNode: newNode(ROOT),
	}

	zipReader, err = zip.OpenReader(wordsFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read zip file: %w", err)
	}
	defer zipReader.Close()

	for _, zipfile := range zipReader.File {
		var (
			file    io.ReadCloser
			scanner *bufio.Scanner
			err     error
		)

		file, err = zipfile.Open()
		if err != nil {
			return nil, fmt.Errorf("could not read file content: %w", err)
		}

		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.ToLower(strings.ReplaceAll(scanner.Text(), " ", ""))
			trie.insert(word)
		}
	}

	return trie, nil
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
