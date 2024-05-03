package minienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateFile(t *testing.T, filename string, lines []string) {
	file, err := os.Create(filename)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	for _, l := range lines {
		_, err := file.WriteString(l + "\n")
		if err != nil {
			assert.FailNow(t, err.Error())
		}
	}
}

func RemoveFile(t *testing.T, filename string) {
	if err := os.Remove(filename); err != nil {
		assert.FailNow(t, err.Error())
	}
}
