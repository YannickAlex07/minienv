package tag

import (
	"errors"
	"fmt"
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

func ParseMinienvTag(tagStr string) (MinienvTag, error) {
	tokens := tokenize(tagStr)
	if len(tokens) == 0 {
		return MinienvTag{}, &ParsingError{
			TagString: tagStr,
			Msg:       "tag is empty",
		}
	}

	p := newParser(tokens)
	parsed, err := p.Parse()
	if err != nil {
		return MinienvTag{}, &ParsingError{
			TagString: tagStr,
			Msg:       err.Error(),
		}
	}

	return parsed, nil
}

// TOKENIZATION

// Tokenize takes a minienv tag string and splits it into the different valid tokens.
// To keep the tokenizations simple, the valid tokens are currently:
//   - ,
//   - =
//   - [ and ]
//   - any literal
//
// Fully empty tokens are just ignored, however spaces in literal are allowed
func tokenize(tag string) []string {
	token := ""
	tokens := []string{}

	for _, char := range tag {
		switch char {
		case ',', '=', '[', ']':
			trimmed := strings.TrimSpace(token)
			if trimmed != "" {
				tokens = append(tokens, trimmed)
				token = ""
			}

			tokens = append(tokens, string(char))
		default:
			token += string(char)
		}
	}

	trimmed := strings.TrimSpace(token)
	if trimmed != "" {
		tokens = append(tokens, trimmed)
	}

	return tokens
}

// PARSER

// This parser is used to parse a set of tokens into the different options that minienv supports.
// The need for a custom parser comes from the fact that the tag does not only allow commas for
// seperation but also for splitting values and within defaults. A simple split on a comma would therefore
// not work anymore and a custom parser is required to fully handle all necessary cases that can occur.
type parser struct {
	tokens       []string
	currentIndex int
}

func newParser(tokens []string) *parser {
	return &parser{
		tokens:       tokens,
		currentIndex: 0,
	}
}

// Gets the current token or an empty string if we are at the end of the token list
func (p *parser) currentToken() string {
	if p.currentIndex >= len(p.tokens) {
		return ""
	}

	return p.tokens[p.currentIndex]
}

// Gets the next token or an empty string if it is at the end of the token list
func (p *parser) peek() string {
	if p.currentIndex+1 >= len(p.tokens) {
		return ""
	}

	return p.tokens[p.currentIndex+1]
}

// Parses the split option. The split option is of the form "split=<token>".
func (p *parser) parseSplitOption() (string, error) {
	// consume the split token
	p.currentIndex++

	// now we should see an = sign
	if ct := p.currentToken(); ct != "=" {
		return "", errors.New("split tag is missing = sign")
	}
	p.currentIndex++

	// if after the current value there is not a comma or EOF, the split tag itself is invalid
	if next := p.peek(); next != "" && next != "," {
		return "", errors.New("invalid number of tokens in split tag")
	}

	// check if the current token is actually present
	splitToken := p.currentToken()
	if splitToken == "" {
		return "", errors.New("missing split token")
	}

	return splitToken, nil
}

// Parses the default option. The default option is of the form "default=<value>".
// To allow for splittable values, the default can be wrapped in []. The value inside
// the [] can use characters like ",", which otherwise would cause parsing issues.
func (p *parser) parseDefaultOption() (string, error) {
	// consume the default token
	p.currentIndex++

	// now we should see an = sign
	if ct := p.currentToken(); ct != "=" {
		return "", errors.New("default tag is missing = sign")
	}
	p.currentIndex++

	// the token now will be our start token
	start := p.currentToken()
	switch start {
	case "", ",":
		return "", fmt.Errorf("invalid default token: %s", start)

	case "[":
		// if we see a [ we need to consume until we see a ]
		startingIndex := p.currentIndex
		for {
			// if we run out of token before seeing a ], it is invalid
			if p.currentIndex >= len(p.tokens) {
				return "", errors.New("missing closing ] in default tag")
			}

			// we found the closeing ], therefore we can break
			if p.currentToken() == "]" {
				break
			}

			// otherwise it is included in the default
			p.currentIndex++
		}

		// we will ingore the opening [ and closing ] and return the value
		val := strings.Join(p.tokens[startingIndex+1:p.currentIndex], "")
		return val, nil
	default:
		// if we see a valid token that is neither a comma nor a [, it is the default value
		return start, nil
	}
}

// Parses all tokens into the different options that minienv supports.
// Currently this includes:
//   - optional = if the field is optional
//   - split = if the field is a slice, this is the token that we use to split
//   - default = the default value for the field
//
// This custom parser is needed because the "," character is not only used to seperate
// the different options, but can also be used as a splittable character or within the
// default value.
func (p *parser) Parse() (MinienvTag, error) {
	tag := MinienvTag{}

	for p.currentIndex < len(p.tokens) {
		token := p.tokens[p.currentIndex]

		// the first token is always the lookup name
		if p.currentIndex == 0 {
			tag.LookupName = p.tokens[p.currentIndex]
			p.currentIndex++
			continue
		}

		switch token {
		case "optional":
			tag.Optional = true

		case "split":
			val, err := p.parseSplitOption()
			if err != nil {
				return MinienvTag{}, err
			}

			tag.SplitOn = val

		case "default":
			def, err := p.parseDefaultOption()
			if err != nil {
				return MinienvTag{}, err
			}

			tag.Default = def

		case ",":
			// skip commas as they are just seperators
			p.currentIndex++
			continue

		default:
			// if we don't recognize the token, something is invalid
			return MinienvTag{}, fmt.Errorf("invalid token in tag \"%s\"", token)
		}

		p.currentIndex++ // move to the next token
	}

	return tag, nil
}
