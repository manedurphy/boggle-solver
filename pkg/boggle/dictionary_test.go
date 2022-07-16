package boggle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDictionary(t *testing.T) {
	tests := []struct {
		name          string
		wordsFilePath string
		err           bool
	}{
		{
			name: "Fail to read zip",
			err:  true,
		},
		{
			name:          "Successfully initialize dictionary",
			wordsFilePath: "../../data/words.zip",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := InitDictionary(test.wordsFilePath)
			if test.err {
				assert.NotNil(t, err)
				return
			}

			assert.Nil(t, err)
		})
	}
}
