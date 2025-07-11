package minienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createFile(t *testing.T, filename string, lines []string) {
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

func removeFile(t *testing.T, filename string) {
	if err := os.Remove(filename); err != nil {
		assert.FailNow(t, err.Error())
	}
}

func setenv(t *testing.T, key, value string) {
	t.Helper()

	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable %s: %v", key, err)
	}

	t.Cleanup(func() {
		err := os.Unsetenv(key)
		if err != nil {
			t.Errorf("Failed to unset environment variable %s: %v", key, err)
		}
	})
}
