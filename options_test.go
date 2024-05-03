package minienv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv"
)

func TestWithOverrides(t *testing.T) {
	// Arrange
	type S struct {
		FromOverride string `env:"FROM_OVERRIDE"`
		FromEnv      string `env:"FROM_ENV"`
	}

	os.Setenv("FROM_OVERRIDE", "from-env")
	defer os.Unsetenv("FROM_OVERRIDE")

	os.Setenv("FROM_ENV", "from-env")
	defer os.Unsetenv("FROM_ENV")

	overrides := map[string]string{
		"FROM_OVERRIDE": "from-override",
	}

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithOverrides(overrides))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "from-override", s.FromOverride)
	assert.Equal(t, "from-env", s.FromEnv)
}

func TestWithEnvFile(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithFile(filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "value", s.Value)
}

func TestWithQuotedEnvFile(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithFile(filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "double", s.Double)
	assert.Equal(t, "single", s.Single)
}

func TestWithEnvFileAndMissingFile(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	os.Setenv("VALUE", "val")
	defer os.Unsetenv("VALUE")

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithFile(filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithRequiredEnvFileAndMissingFile(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"VALUE"`
	}

	filename := "test.env" // file does not exist

	// Act
	var s S
	err := minienv.Load(&s, minienv.WithRequiredFile(filename))

	// Assert
	assert.NotNil(t, err)
}

func TestWithRequiredEnvFile(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithRequiredFile(filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithFileAndDefaultFile(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithFile())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}

func TestWithMultipleFiles(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithFile(filename1, filename2))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "one", s.One)
	assert.Equal(t, "two", s.Two)
}

func TestWithEmptyLines(t *testing.T) {
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
	err := minienv.Load(&s, minienv.WithFile(filename))

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "val", s.Value)
}
