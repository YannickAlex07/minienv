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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, "required", s.Req)
	assert.Equal(t, "", s.Opt) // should be empty as it was never set
	assert.Equal(t, "optionalexists", s.OptEx)
}

func TestLoadWithNonPointer(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE"`
	}

	// Act
	var s S
	err := minienv.Load(s)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, minienv.ErrInvalidInput, err)
}

func TestLoadWithNonStruct(t *testing.T) {
	// Act
	var s string
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, minienv.ErrInvalidInput, err)
}

func TestLoadWithMixedTags(t *testing.T) {
	// Arrange
	type S struct {
		Value     string `env:"TEST_VALUE"`
		NotTagged string
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "test-value", s.Value)
	assert.Equal(t, "", s.NotTagged)
}

func TestLoadWithMissingValue(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	missingErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", missingErr.Field)
	assert.ErrorContains(t, missingErr, "required field has no value and no default")
}

func TestLoadWithUnsupportedType(t *testing.T) {
	// Arrange
	type S struct {
		Value map[string]string `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.CoversionError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.Equal(t, "test-value", conversionErr.Value)
	assert.ErrorContains(t, conversionErr, "unsupported type")
}

func TestLoadWithInvalidInt(t *testing.T) {
	// Arrange
	type S struct {
		Value int `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.CoversionError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.Equal(t, "test-value", conversionErr.Value)
	assert.ErrorContains(t, conversionErr, "strconv.Atoi: parsing \"test-value\": invalid syntax")
}

func TestLoadWithInvalidBool(t *testing.T) {
	// Arrange
	type S struct {
		Value bool `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.CoversionError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.Equal(t, "test-value", conversionErr.Value)
	assert.ErrorContains(t, conversionErr, "parsing \"test-value\": invalid syntax")
}

func TestLoadWithInvalidFloat(t *testing.T) {
	// Arrange
	type S struct {
		Value float64 `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.CoversionError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.Equal(t, "test-value", conversionErr.Value)
	assert.ErrorContains(t, conversionErr, "parsing \"test-value\": invalid syntax")
}

func TestLoadWithDefaultString(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE,default=Hello"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "Hello", s.Value)
}

func TestLoadWithDefaultInt(t *testing.T) {
	// Arrange
	type S struct {
		Value int `env:"TEST_VALUE,default=5"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 5, s.Value)
}

func TestLoadWithDefaultMissingValue(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE,default"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	tagParseErr := err.(minienv.TagParsingError)
	assert.Equal(t, "Value", tagParseErr.Field)
	assert.ErrorContains(t, tagParseErr, "default tag does not contain a single value")
}
