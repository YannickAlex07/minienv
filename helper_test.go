package minienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateEnvFile(t *testing.T, filename string, vars map[string]string) {
	file, err := os.Create(filename)
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	for k, v := range vars {
		_, err := file.WriteString(k + "=" + v + "\n")
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
