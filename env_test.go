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

	missingErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", missingErr.Field)
	assert.ErrorContains(t, missingErr, "no value was found for required field with lookup key TEST_VALUE")
}

func TestLoadWithMissingNestedValue(t *testing.T) {
	// Arrange
	type S struct {
		N struct {
			Value string `env:"TEST_VALUE"`
		}
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	missingErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", missingErr.Field)
	assert.ErrorContains(t, missingErr, "no value was found for required field with lookup key TEST_VALUE")
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

	conversionErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", conversionErr.Field)
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

	conversionErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", conversionErr.Field)
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

	conversionErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", conversionErr.Field)
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

	conversionErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", conversionErr.Field)
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

	tagParseErr := err.(minienv.LoadError)
	assert.Equal(t, "Value", tagParseErr.Field)
	assert.ErrorContains(t, tagParseErr, "default tag is missing = sign")
}

func TestLoadWithUnsettableField(t *testing.T) {
	// Arrange
	type S struct {
		_ string `env:"TEST_VALUE"`
	}

	os.Setenv("TEST_VALUE", "test-value")
	defer os.Unsetenv("TEST_VALUE")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not valid or cannot be set")
}

func TestLoadWithSplittableFloatField(t *testing.T) {
	type S struct {
		Floats        []float64 `env:"TEST_FLOATS,split=,"`
		FloatDefaults []float64 `env:"TEST_FLOATS_DEF,split=,,default=[1.1,2.2,3.3]"`
	}

	os.Setenv("TEST_FLOATS", "1.1,2.2,3.3")
	defer os.Unsetenv("TEST_FLOATS")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.Floats)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.FloatDefaults)
}

func TestLoadWithSplittableStringField(t *testing.T) {
	type S struct {
		Str        []string `env:"TEST_STR,split=,"`
		StrDefault []string `env:"TEST_STR_DEF,split=,,default=[test1,test2]"`
	}

	os.Setenv("TEST_STR", "test1,test2")
	defer os.Unsetenv("TEST_STR")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"test1", "test2"}, s.Str)
	assert.Equal(t, []string{"test1", "test2"}, s.StrDefault)
}

func TestLoadWithSplittableIntField(t *testing.T) {
	type S struct {
		Numbers        []int `env:"TEST_NUMBERS,split=,"`
		NumbersDefault []int `env:"TEST_NUMBERS_DEF,split=,,default=[1,2,3]"`
	}

	os.Setenv("TEST_NUMBERS", "1,2,3")
	defer os.Unsetenv("TEST_NUMBERS")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, s.Numbers)
	assert.Equal(t, []int{1, 2, 3}, s.NumbersDefault)
}

func TestLoadWithSplittableBoolField(t *testing.T) {
	type S struct {
		Bools        []bool `env:"TEST_BOOLS,split=,"`
		BoolsDefault []bool `env:"TEST_BOOLS_DEF,split=,,default=[true,false]"`
	}

	os.Setenv("TEST_BOOLS", "true,false")
	defer os.Unsetenv("TEST_BOOLS")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false}, s.Bools)
	assert.Equal(t, []bool{true, false}, s.BoolsDefault)
}

func TestLoadWithSplittableUnsupportedType(t *testing.T) {
	type S struct {
		Unsupported []struct{} `env:"TEST_UNSUPPORTED,split=,"`
	}

	os.Setenv("TEST_UNSUPPORTED", "test1,test2")
	defer os.Unsetenv("TEST_UNSUPPORTED")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.ErrorContains(t, err, "failed to convert value test1 in slice to type struct")
}
