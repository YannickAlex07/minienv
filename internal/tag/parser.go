package tag

import (
	"fmt"
	"regexp"
	"strings"
)

// ERROR

type ParsingError struct {
	TagString string
	Msg       string
}

func (e *ParsingError) Error() string {
	return fmt.Sprintf("failed to parse tag \"%s\": %s", e.TagString, e.Msg)
}

// TAG

type MinienvTag struct {
	// The name that will be used when looking up the environment variable
	LookupName string

	// If this is set to true, the field is optional
	Optional bool

	// This represents the default value for the field in case the env variable is not set
	Default string

	// In case the field is a slice, this is the string that we use to split the value on
	SplitOn string
}

// PARSER

// Compile the regex pattern that is used to parse the tag.
// We will use MustCompile here, as we are sure that this will work.
func compilePattern() *regexp.Regexp {
	patterns := []string{
		`(?P<optional>optional)`,                         // optional
		`default=(?P<default>[\w\d\s]+|\[[\w\W\d\s]+\])`, // default=<token> or default=[<token>, <token>, ...]
		`split=(?P<split>[\w\W\d])`,                      // split=<token>
	}

	// the final pattern will essentially be: <lookup>(,optional|default=<token>|split=<token>)*
	opStr := strings.Join(patterns, "|")
	pattern := fmt.Sprintf(`^(?P<lookup>[a-zA-Z_-]+)(?:,(?:%s))*$`, opStr)

	return regexp.MustCompile(pattern)
}

// we compile it once gloabally, so we don't have to do it every time we parse a tag
var pattern *regexp.Regexp = compilePattern()

func Parse(tagStr string) (MinienvTag, error) {
	if tagStr == "" {
		return MinienvTag{}, &ParsingError{
			TagString: tagStr,
			Msg:       "tag is empty",
		}
	}

	matches := pattern.FindStringSubmatch(tagStr)
	if matches == nil {
		return MinienvTag{}, &ParsingError{
			TagString: tagStr,
			Msg:       "invalid tag format",
		}
	}

	options := map[string]string{}
	for i, name := range pattern.SubexpNames() {
		if name == "" {
			continue
		}

		options[name] = matches[i]
	}

	tag := MinienvTag{
		LookupName: options["lookup"],
		Optional:   options["optional"] != "",
		Default:    strings.Trim(options["default"], "[]"),
		SplitOn:    options["split"],
	}

	return tag, nil
}
