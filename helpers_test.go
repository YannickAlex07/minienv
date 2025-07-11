package minienv_test

import (
	"os"
	"testing"
)

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
