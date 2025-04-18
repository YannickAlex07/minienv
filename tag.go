package minienv

import (
	"errors"
	"reflect"
	"strings"
)

// This struct hold all the metadata about a found "env"-tag for a field
type tag struct {
	// This is the name of the env variable we need to look for
	name string

	// This is a flag that tells us if the variable is required
	required bool

	// This is the default value for the variable, can be empty and therefore invalid
	defaultValue string

	// This is the character that we use to split the value in case of slices
	splitOnChar rune
}

// Parses the `env` tag and returns the bundled information about the tag.
// The first return value is the tag itself, the second return value is a flag indicating if the tag was found
// and the third return value is an error if the tag was invalid.
func parseTag(field reflect.StructField) (tag, bool, error) {
	required := true
	var defaultVal string

	value, found := field.Tag.Lookup("env")
	if !found {
		return tag{}, false, nil
	}

	// check any tag options
	parts := strings.Split(value, ",")
	for _, p := range parts[1:] {
		trimmed := strings.TrimSpace(p)
		splitted := strings.Split(trimmed, "=")

		// tag is optional
		switch splitted[0] {
		case "optional":
			required = false

		case "default":

			// if we have more or less than 2 elements we have an invalid tag
			if len(splitted) != 2 {
				return tag{}, true, errors.New("invalid default tag")
			}

			defaultVal = splitted[1]
		}
	}

	return tag{
		name:         parts[0],
		required:     required,
		defaultValue: defaultVal,
	}, true, nil
}
