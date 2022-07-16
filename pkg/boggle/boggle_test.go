package boggle

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	WORDS_FILE_PATH = "../../data/words.zip"
)

func init() {
	err := InitDictionary(WORDS_FILE_PATH)
	if err != nil {
		panic(err)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		board        [][]string
	}{
		{
			name:         "Fails for invalid boggle board",
			errorMessage: "board is invalid",
			board:        [][]string{},
		},
		{
			name:         "Fails for imbalanced boggle board",
			errorMessage: "board is imbalanced",
			board: [][]string{
				{"A", "B", "C"},
				{"D", "E"},
			},
		},
		{
			name:         "Successfully creates new boggle instance",
			errorMessage: "board is imbalanced",
			board: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boggle, err := New(test.board)
			if err != nil {
				assert.Equal(t, errors.New(test.errorMessage), err)
				return
			}

			assert.NotNil(t, boggle)
		})
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name   string
		board  [][]string
		result *Result
	}{
		{
			name: "Solve 1x3 board",
			board: [][]string{
				{"A", "B", "C"},
			},
			result: &Result{
				WordsFound: []string{"abc"},
				Score:      1,
			},
		},
		{
			name: "Solve 2x3 board",
			board: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
			},
			result: &Result{
				WordsFound: []string{
					"abc",
					"abd",
					"abe",
					"abed",
					"ade",
					"bad",
					"bade",
					"bae",
					"bcf",
					"bde",
					"bea",
					"bead",
					"bec",
					"bed",
					"bef",
					"dab",
					"dae",
					"dea",
					"deb",
					"dec",
					"def",
					"ead",
					"ecb",
					"fec",
					"fed",
				},
				Score: 25,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			boggle, err := New(test.board)
			if err != nil {
				t.Errorf("could not create boggle board: %s", err)
				return
			}

			result := boggle.Solve()
			assert.Equal(t, test.result, result)
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	tests := []struct {
		name  string
		board [][]string
	}{
		{
			name: "Solve 1x3 board",
			board: [][]string{
				{"A", "B", "C"},
			},
		},
		{
			name: "Solve 2x3 board",
			board: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
			},
		},
		{
			name: "Solve 3x3 board",
			board: [][]string{
				{"A", "B", "C"},
				{"D", "E", "F"},
				{"H", "I", "J"},
			},
		},
		{
			name: "Solve 4x4 board",
			board: [][]string{
				{"D", "A", "T", "T"},
				{"U", "A", "K", "Q"},
				{"P", "L", "A", "U"},
				{"M", "Y", "O", "P"},
				{"O", "G", "G", "M"},
				{"L", "A", "N", "A"},
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				boggle, err := New(test.board)
				if err != nil {
					b.Errorf("could not create boggle board: %s", err)
					return
				}

				boggle.Solve()
			}
		})
	}
}
