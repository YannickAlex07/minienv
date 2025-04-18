package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv/internal/tag"
)

func TestParseEnvTagWithCompleteValidOptions(t *testing.T) {
	// Arrange
	variations := []string{
		"TEST,split=,,default=[10,20,30],optional",
		"TEST,default=[10,20,30],split=,,optional",
		"TEST,optional,default=[10,20,30],split=,",
	}

	expected := tag.MinienvTag{
		LookupName: "TEST",
		Optional:   true,
		Default:    "10,20,30",
		SplitOn:    ",",
	}

	// Act
	for _, variation := range variations {
		t.Run(variation, func(t *testing.T) {
			t.Parallel()

			result, err := tag.ParseMinienvTag(variation)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	}
}

// TESTS RELATED TO THE SPLIT OPTION

func TestParseEnvTagWithValidSplitOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name     string
		tagStr   string
		expected tag.MinienvTag
	}

	testCases := []testCase{
		{
			name:   "Split with no other options",
			tagStr: "TEST,split=,",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ",",
			},
		},
		{
			name:   "Split on colon with no other options",
			tagStr: "TEST,split=:",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ":",
			},
		},
		{
			name:   "Split on comma with other options",
			tagStr: "TEST,split=,,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ",",
				Optional:   true,
			},
		},
		{
			name:   "Split on colon with other options",
			tagStr: "TEST,split=:,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ":",
				Optional:   true,
			},
		},
		{
			name:   "Split on comma as the last option",
			tagStr: "TEST,optional,split=,",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ",",
				Optional:   true,
			},
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := tag.ParseMinienvTag(tCase.tagStr)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tCase.expected, result)
		})
	}
}

func TestParseEnvTagWithInvalidSplitOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name   string
		tagStr string
	}

	testCases := []testCase{
		{
			name:   "Split with missing = and missing character",
			tagStr: "TEST,split",
		},
		{
			name:   "Split with missing character",
			tagStr: "TEST,split=",
		},
		{
			name:   "Split with missing character and no other options",
			tagStr: "TEST,split=,optional",
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := tag.ParseMinienvTag(tCase.tagStr)

			// Assert
			assert.Error(t, err)
		})
	}
}

// TESTS RELATED TO THE DEFAULT OPTION
func TestParseEnvTagWithValidDefaultOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name     string
		tagStr   string
		expected tag.MinienvTag
	}

	testCases := []testCase{
		{
			name:   "Default with a simple value",
			tagStr: "TEST,default=something",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something",
			},
		},
		{
			name:   "Default with a simple value that contains spaces",
			tagStr: "TEST,default=something interesting",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something interesting",
			},
		},
		{
			name:   "Default with a slice",
			tagStr: "TEST,default=[10,20,30]",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "10,20,30",
			},
		},
		{
			name:   "Default with a simple value and other options",
			tagStr: "TEST,default=something,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something",
				Optional:   true,
			},
		},
		{
			name:   "Default with a slice and other options",
			tagStr: "TEST,default=[10,20,30],optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "10,20,30",
				Optional:   true,
			},
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := tag.ParseMinienvTag(tCase.tagStr)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tCase.expected, result)
		})
	}
}

func TestParseEnvTagWithInvalidDefaultOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name   string
		tagStr string
	}

	testCases := []testCase{
		{
			name:   "Default with missing =",
			tagStr: "TEST,default",
		},
		{
			name:   "Default with missing value",
			tagStr: "TEST,default=",
		},
		{
			name:   "Default with non-closed slice",
			tagStr: "TEST,default=[10,20,",
		},
		{
			name:   "Default with non-closed slice and other options",
			tagStr: "TEST,default=[10,20,optional",
		},
		{
			name:   "Default with missing value and other options",
			tagStr: "TEST,default=,optional",
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := tag.ParseMinienvTag(tCase.tagStr)

			// Assert
			assert.Error(t, err)
		})
	}
}
