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
