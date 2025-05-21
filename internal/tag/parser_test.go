package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv/internal/tag"
)

func TestParseEnvTagWithCompleteValidOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name     string
		tagStr   string
		expected tag.MinienvTag
	}

	testCases := []testCase{
		{
			name:   "complete tag with all options",
			tagStr: "TEST,split=,,default=[10,20,30],optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Optional:   true,
				Default:    "10,20,30",
				SplitOn:    ",",
			},
		},
		{
			name:   "complete tag with all options",
			tagStr: "TEST,default=[10,20,30],split=,,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Optional:   true,
				Default:    "10,20,30",
				SplitOn:    ",",
			},
		},
		{
			name:   "complete tag with all options",
			tagStr: "TEST,optional,default=[10,20,30],split=,",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Optional:   true,
				Default:    "10,20,30",
				SplitOn:    ",",
			},
		},
		{
			name:   "just the lookup name",
			tagStr: "TEST",
			expected: tag.MinienvTag{
				LookupName: "TEST",
			},
		},
	}

	// Act
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result, err := tag.Parse(testCase.tagStr)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestParseEnvTagWithInvalidOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name          string
		tagStr        string
		errorContains string
	}

	testCases := []testCase{
		{
			name:          "empty tag",
			tagStr:        "",
			errorContains: "tag is empty",
		},
		{
			name:          "unknown option",
			tagStr:        "TEST,unknown",
			errorContains: "invalid tag format",
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := tag.Parse(tCase.tagStr)

			// Assert
			assert.ErrorContains(t, err, tCase.errorContains)
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
			name:   "split with no other options",
			tagStr: "TEST,split=,",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ",",
			},
		},
		{
			name:   "split on colon with no other options",
			tagStr: "TEST,split=:",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ":",
			},
		},
		{
			name:   "split on comma with other options",
			tagStr: "TEST,split=,,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ",",
				Optional:   true,
			},
		},
		{
			name:   "split on colon with other options",
			tagStr: "TEST,split=:,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				SplitOn:    ":",
				Optional:   true,
			},
		},
		{
			name:   "split on comma as the last option",
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

			result, err := tag.Parse(tCase.tagStr)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tCase.expected, result)
		})
	}
}

func TestParseEnvTagWithInvalidSplitOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name          string
		tagStr        string
		errorContains string
	}

	testCases := []testCase{
		{
			name:          "split with missing = and missing character",
			tagStr:        "TEST,split",
			errorContains: "invalid tag format",
		},
		{
			name:          "split with missing character",
			tagStr:        "TEST,split=",
			errorContains: "invalid tag format",
		},
		{
			name:          "split with missing character and no other options",
			tagStr:        "TEST,split=,optional",
			errorContains: "invalid tag format",
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := tag.Parse(tCase.tagStr)

			// Assert
			assert.ErrorContains(t, err, tCase.errorContains)
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
			name:   "default with a simple value",
			tagStr: "TEST,default=something",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something",
			},
		},
		{
			name:   "default with a simple value that contains spaces",
			tagStr: "TEST,default=something interesting",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something interesting",
			},
		},
		{
			name:   "default with a slice",
			tagStr: "TEST,default=[10,20,30]",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "10,20,30",
			},
		},
		{
			name:   "default with a simple value and other options",
			tagStr: "TEST,default=something,optional",
			expected: tag.MinienvTag{
				LookupName: "TEST",
				Default:    "something",
				Optional:   true,
			},
		},
		{
			name:   "default with a slice and other options",
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

			result, err := tag.Parse(tCase.tagStr)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tCase.expected, result)
		})
	}
}

func TestParseEnvTagWithInvalidDefaultOptions(t *testing.T) {
	// Arrange
	type testCase struct {
		name          string
		tagStr        string
		errorContains string
	}

	testCases := []testCase{
		{
			name:          "Default with missing =",
			tagStr:        "TEST,default",
			errorContains: "invalid tag format",
		},
		{
			name:          "Default with missing value",
			tagStr:        "TEST,default=",
			errorContains: "invalid tag format",
		},
		{
			name:          "Default with non-closed slice",
			tagStr:        "TEST,default=[10,20,",
			errorContains: "invalid tag format",
		},
		{
			name:          "Default with non-closed slice and other options",
			tagStr:        "TEST,default=[10,20,optional",
			errorContains: "invalid tag format",
		},
		{
			name:          "Default with missing value and other options",
			tagStr:        "TEST,default=,optional",
			errorContains: "invalid tag format",
		},
	}

	// Act
	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := tag.Parse(tCase.tagStr)

			// Assert
			assert.ErrorContains(t, err, tCase.errorContains)
		})
	}
}
