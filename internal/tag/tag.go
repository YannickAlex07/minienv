package tag

import (
	"errors"
	"strings"
)

type Tag struct {
	// If this is set to true, the field is optional
	Optional bool

	// This represents the default value for the field in case the env variable is not set
	Default string

	// In case the field is a slice, this is the character that we use to split the value
	SplitOn rune
}

func ParseEnvTag(tagStr string) (Tag, error) {
	parser, err := newTagParser(tagStr)
	if err != nil {
		return Tag{}, err
	}

	tag := Tag{}
	for token, err := range parser.Parse() {
		if err != nil {
			return Tag{}, err
		}

		parts := strings.Split(token, "=")
		switch parts[0] {
		case "optional":
			tag.Optional = true

		case "default":
			tag.Default = parts[1]

		case "split":
			tag.SplitOn = rune(parts[1][0])

		default:
			return Tag{}, errors.New("unknown tag option: " + parts[0])
		}
	}

	return tag, nil
}
