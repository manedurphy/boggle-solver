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

var (
	dictionary *dictionaryRoot
)

type (
	node struct {
		char     string
		isWord   bool
		children [SIZE]*node
	}

	dictionaryRoot struct {
		RootNode *node
	}
)

func newNode(char string) *node {
	return &node{
		char:   char,
		isWord: false,
	}
}

func InitDictionary(wordsFilePath string) error {
	var (
		dRoot     *dictionaryRoot
		zipReader *zip.ReadCloser
		err       error
	)

	dRoot = &dictionaryRoot{
		RootNode: newNode(ROOT),
	}

	zipReader, err = zip.OpenReader(wordsFilePath)
	if err != nil {
		return fmt.Errorf("could not read zip file: %w", err)
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
			return fmt.Errorf("could not read file content: %w", err)
		}

		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.ToLower(strings.ReplaceAll(scanner.Text(), " ", ""))
			dRoot.insert(word)
		}
	}

	dictionary = dRoot
	return nil
}

func (d *dictionaryRoot) insert(word string) {
	current := d.RootNode
	strippedWord := strings.ToLower(strings.ReplaceAll(word, " ", ""))

	for i := 0; i < len(strippedWord); i++ {
		char := string(strippedWord[i])
		index := strippedWord[i] - 'a'

		if current.children[index] == nil {
			current.children[index] = newNode(char)
		}

		current = current.children[index]
	}
	current.isWord = true
}
