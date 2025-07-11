package minienv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv"
)

func TestWithFallbackValues(t *testing.T) {
	t.Parallel()

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

func TestWithFile(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"FROM_FILE"`
	}

	// create env file
	filename := "test.env"

	CreateFile(t, filename, []string{
		"FROM_FILE=value",
	})
	defer RemoveFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false, filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "value", s.Value)
}

func TestWithFileAndQuoted(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Double string `env:"DOUBLE"`
		Single string `env:"SINGLE"`
	}

	// create env file
	filename := "test.env"

	CreateFile(t, filename, []string{
		"DOUBLE=\"double\"",
		"SINGLE='single'",
	})
	defer RemoveFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false, filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "double", s.Double)
	assert.Equal(t, "single", s.Single)
}

func TestWithFileAndMissingOptionalFile(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	setenv(t, "VALUE", "val")

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false, filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithFileAndMissingRequiredFile(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(true, filename))

	// Assert
	assert.NotNil(t, err)
}

func TestWithFileAndRequired(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env"

	CreateFile(t, filename, []string{
		"VALUE=val",
	})
	defer RemoveFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(true, filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithFileAndDefaultFile(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := ".env"

	CreateFile(t, filename, []string{
		"VALUE=val",
	})
	defer RemoveFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithMultipleFiles(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		One string `env:"ONE"`
		Two string `env:"TWO"`
	}

	filename1 := "one.env"
	filename2 := "two.env"

	CreateFile(t, filename1, []string{
		"ONE=one",
	})
	defer RemoveFile(t, filename1)

	CreateFile(t, filename2, []string{
		"TWO=two",
	})
	defer RemoveFile(t, filename2)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false, filename1, filename2))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "one", s.One)
	assert.Equal(t, "two", s.Two)
}

func TestWithEmptyLines(t *testing.T) {
	t.Parallel()

	// Arrange
	type S struct {
		Value string `env:"VAL"`
	}

	filename := "test.env"

	CreateFile(t, filename, []string{
		"VAL=val",
		"",
		"# comment",
	})
	defer RemoveFile(t, filename)

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(false, filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithPrefix(t *testing.T) {
	t.Parallel()

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
