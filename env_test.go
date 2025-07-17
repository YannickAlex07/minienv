package minienv_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv"
)

func TestLoadWithString(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-string")

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

	setenv(t, "TEST_VALUE", "3823992")

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

	setenv(t, "TEST_VALUE", "34.3243")

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

	setenv(t, "TEST_VALUE", "true")

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

	setenv(t, "TEST_VALUE", "test")

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

	setenv(t, "REQ", "required")
	setenv(t, "OPT_EX", "optionalexists")

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

	setenv(t, "TEST_VALUE", "test-value")

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
	assert.ErrorContains(t, missingErr, "no value was found for field with lookup key: TEST_VALUE")
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

	missingErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", missingErr.Field)
	assert.ErrorContains(t, missingErr, "no value was found for field with lookup key: TEST_VALUE")
}

func TestLoadWithUnsupportedType(t *testing.T) {
	// Arrange
	type S struct {
		Value complex128 `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.ErrorContains(t, conversionErr, "unsupported type")
}

func TestLoadWithEmptyTag(t *testing.T) {
	// Arrange
	type S struct {
		Value map[string]string `env:""`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.ErrorContains(t, conversionErr, "tag string cannot be empty")
}

func TestLoadWithUnknownTagOption(t *testing.T) {
	// Arrange
	type S struct {
		Value map[string]string `env:"TEST_VALUE,unknown=option"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.ErrorContains(t, conversionErr, "unknown tag option \"unknown\"")
}

func TestLoadWithInvalidInt(t *testing.T) {
	// Arrange
	type S struct {
		Value int `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.ErrorContains(t, conversionErr, "strconv.Atoi: parsing \"test-value\": invalid syntax")
}

func TestLoadWithInvalidBool(t *testing.T) {
	// Arrange
	type S struct {
		Value bool `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
	assert.Equal(t, "Value", conversionErr.Field)
	assert.ErrorContains(t, conversionErr, "parsing \"test-value\": invalid syntax")
}

func TestLoadWithInvalidFloat(t *testing.T) {
	// Arrange
	type S struct {
		Value float64 `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)

	conversionErr := err.(minienv.FieldError)
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

func TestLoadWithDefaultStringIncludingEqualSign(t *testing.T) {
	// Arrange
	type S struct {
		Value string `env:"TEST_VALUE,default=key=value"`
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "key=value", s.Value)
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

	parseError := err.(minienv.FieldError)
	assert.Equal(t, "Value", parseError.Field)
	assert.ErrorContains(t, parseError, "default env value cannot be empty")
}

func TestLoadWithUnsettableField(t *testing.T) {
	// Arrange
	type S struct {
		_ string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test-value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not valid or cannot be set")
}

func TestLoadSliceWithFloat(t *testing.T) {
	type S struct {
		Floats        []float64 `env:"TEST_FLOATS"`
		FloatDefaults []float64 `env:"TEST_FLOATS_DEF,default=1.1|2.2|3.3"`
	}

	setenv(t, "TEST_FLOATS", "1.1|2.2|3.3")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.Floats)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, s.FloatDefaults)
}

func TestLoadSliceWithString(t *testing.T) {
	type S struct {
		Str        []string `env:"TEST_STR"`
		StrDefault []string `env:"TEST_STR_DEF,default=test1|test2"`
	}

	setenv(t, "TEST_STR", "test1|test2")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []string{"test1", "test2"}, s.Str)
	assert.Equal(t, []string{"test1", "test2"}, s.StrDefault)
}

func TestLoadSliceWithInt(t *testing.T) {
	type S struct {
		Numbers        []int `env:"TEST_NUMBERS"`
		NumbersDefault []int `env:"TEST_NUMBERS_DEF,default=1|2|3"`
	}

	setenv(t, "TEST_NUMBERS", "1|2|3")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, s.Numbers)
	assert.Equal(t, []int{1, 2, 3}, s.NumbersDefault)
}

func TestLoadSliceWithBool(t *testing.T) {
	type S struct {
		Bools        []bool `env:"TEST_BOOLS"`
		BoolsDefault []bool `env:"TEST_BOOLS_DEF,default=true|false"`
	}

	setenv(t, "TEST_BOOLS", "true|false")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false}, s.Bools)
	assert.Equal(t, []bool{true, false}, s.BoolsDefault)
}

func TestLoadSliceWithUnsupportedType(t *testing.T) {
	type S struct {
		Unsupported []struct{} `env:"TEST_UNSUPPORTED"`
	}

	setenv(t, "TEST_UNSUPPORTED", "test1|test2")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.ErrorContains(t, err, "failed to set slice element 0: unsupported type: struct")
}

func TestLoadMapWithString(t *testing.T) {
	type S struct {
		MapField map[string]string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test1:value1|test2:value2")

	expected := map[string]string{
		"test1": "value1",
		"test2": "value2",
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, s.MapField)
}

func TestLoadMapWithInt(t *testing.T) {
	type S struct {
		MapField map[int]int `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "1:1|2:2")

	expected := map[int]int{
		1: 1,
		2: 2,
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, s.MapField)
}

func TestLoadMapWithFloat(t *testing.T) {
	type S struct {
		MapField map[float64]float64 `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "1.1:1.1|2.2:2.2")

	expected := map[float64]float64{
		1.1: 1.1,
		2.2: 2.2,
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, s.MapField)
}

func TestLoadMapWithBool(t *testing.T) {
	type S struct {
		MapField map[bool]bool `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "true:true|false:false")

	expected := map[bool]bool{
		true:  true,
		false: false,
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, s.MapField)
}

func TestLoadMapWithUnsupportedType(t *testing.T) {
	type S struct {
		MapField map[string]struct{} `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test1:{}|test2:{}")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to set map value for key \"test1\": unsupported type: struct")
}

func TestLoadMapWithDuplicatedKey(t *testing.T) {
	type S struct {
		MapField map[string]string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "test:first|test:second")

	expected := map[string]string{
		"test": "second",
	}

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected, s.MapField)
}

func TestLoadMapWithWrongValueType(t *testing.T) {
	type S struct {
		MapField map[string]int `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "key:value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to set map value for key \"key\": strconv.Atoi: parsing \"value\": invalid syntax")
}

func TestLoadMapWithWrongKeyType(t *testing.T) {
	type S struct {
		MapField map[int]string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "key:value")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to set map key \"key\": strconv.Atoi: parsing \"key\": invalid syntax")
}

func TestLoadMapWithMissingValue(t *testing.T) {
	type S struct {
		MapField map[string]string `env:"TEST_VALUE"`
	}

	setenv(t, "TEST_VALUE", "key")

	// Act
	var s S
	err := minienv.Load(&s)

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "map value must be in the format key:value, got: key")
}
