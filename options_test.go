package minienv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv"
)

func TestWithFallbackValues(t *testing.T) {
	// Arrange
	type S struct {
		FromFallback string `env:"FROM_FALLBACK"`
		FromEnv      string `env:"FROM_ENV"`

		// The env-value takes precedence over the fallback value
		FromBoth string `env:"FROM_BOTH"`
	}

	setenv(t, "FROM_ENV", "from-env")
	setenv(t, "FROM_BOTH", "from-both-env")

	values := map[string]string{
		"FROM_FALLBACK": "from-fallback",
		"FROM_BOTH":     "from-both-fallback",
	}

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFallbackValues(values))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "from-fallback", s.FromFallback)
	assert.Equal(t, "from-env", s.FromEnv)
	assert.Equal(t, "from-both-env", s.FromBoth)
}

func TestWithEnvFile(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"FROM_FILE"`
	}

	// create env file
	filename := "test.env"

	createFile(t, filename, []string{
		"FROM_FILE=value",
	})
	defer removeFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, false))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "value", s.Value)
}

func TestWithEnvFileAndQuoted(t *testing.T) {
	// Arrange
	type S struct {
		Double string `env:"DOUBLE"`
		Single string `env:"SINGLE"`
	}

	// create env file
	filename := "test.env"

	createFile(t, filename, []string{
		"DOUBLE=\"double\"",
		"SINGLE='single'",
	})
	defer removeFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, false))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "double", s.Double)
	assert.Equal(t, "single", s.Single)
}

func TestWithEnvFileAndMissingOptionalFile(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	setenv(t, "VALUE", "val")

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, false))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithEnvFileAndMissingRequiredFile(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, true))

	// Assert
	assert.NotNil(t, err)
}

func TestWithEnvFileAndRequired(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env"

	createFile(t, filename, []string{
		"VALUE=val",
	})
	defer removeFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, true))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithEnvFileAndEmptyLines(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VAL"`
	}

	filename := "test.env"

	createFile(t, filename, []string{
		"VAL=val",
		"",
		"# comment",
	})
	defer removeFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithEnvFile(filename, false))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithPrefix(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	setenv(t, "PREFIX_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithPrefix("PREFIX_"))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "test-value", s.Value)
}
