package minienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv"
)

func TestLoadWithString(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-string")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, "test-string", s.Value)
}

func TestLoadWithInt(t *testing.T) {
	// Arrange
	type S struct {
		Value int `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "3823992")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, 3823992, s.Value)
}

func TestLoadWithFloat(t *testing.T) {
	// Arrange
	type S struct {
		Value float64 `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "34.3243")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, 34.3243, s.Value)
}

func TestLoadWithBool(t *testing.T) {
	// Arrange
	type S struct {
		Value bool `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "true")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, true, s.Value)
}

func TestLoadWithSingleNested(t *testing.T) {
	// Arrange
	type S struct {
		N struct {
			Value string `env:"TEST_VALUE"`
		}
	}

	os.Setenv("TEST_VALUE", "test")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, "test", s.N.Value)
}

func TestLoadWithOptional(t *testing.T) {
	// Arrange
	type S struct {
		Req   string `env:"REQ"`             // is required
		Opt   string `env:"OPT,optional"`    // is optional and not set
		OptEx string `env:"OPT_EX,optional"` // is optional and set
	}

	os.Setenv("REQ", "required")
	defer os.Unsetenv("REQ")

	os.Setenv("OPT_EX", "optionalexists")
	defer os.Unsetenv("OPT_EX")

	// Act
	var s S
	err := minienv.Load(&s)
	if err != nil {
		assert.FailNow(t, "unexpected error: %v", err)
	}

	// Assert
	assert.Equal(t, "required", s.Req)
	assert.Equal(t, "", s.Opt) // should be empty as it was never set
	assert.Equal(t, "optionalexists", s.OptEx)
}
