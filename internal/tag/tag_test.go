package tag_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/minienv/internal/tag"
)

func TestParseEnvTagWithAllValidOptions(t *testing.T) {
	// Arrange
	variations := []string{
		"split=,,default=[10,20,30],optional",
		"default=[10,20,30],split=,,optional",
		"optional,default=[10,20,30],split=,",
	}

	expected := tag.Tag{
		Optional: true,
		Default:  "[10,20,30]",
		SplitOn:  ',',
	}

	// Act
	for _, variation := range variations {
		t.Run(variation, func(t *testing.T) {
			t.Parallel()

			result, err := tag.ParseEnvTag(variation)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, expected, result)
		})
	}
}
