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
		return MinienvTag{}, errors.New("empty tag")
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

func tokenize(tag string) []string {
	token := ""
	tokens := []string{}

	for _, char := range tag {
		switch char {
		case ',', '=', '[', ']':
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}

			tokens = append(tokens, string(char))
		default:
			token += string(char)
		}
	}

	if token != "" {
		tokens = append(tokens, token)
	}

	return tokens
}

// PARSER

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

func (p *parser) currentToken() string {
	if p.currentIndex >= len(p.tokens) {
		return ""
	}

	return p.tokens[p.currentIndex]
}

func (p *parser) peek() string {
	if p.currentIndex+1 >= len(p.tokens) {
		return ""
	}

	return p.tokens[p.currentIndex+1]
}

func (p *parser) parseSplitOption() (string, error) {
	// consume the split token
	p.currentIndex++

	// now we should see an = sign
	if ct := p.currentToken(); ct != "=" {
		return "", errors.New("invalid split tag")
	}
	p.currentIndex++

	// if after the current value there is not a comma or EOF, the split tag itself is invalid
	if next := p.peek(); next != "" && next != "," {
		return "", errors.New("invalid split tag")
	}

	// check if the current token is actually present
	splitToken := p.currentToken()
	if splitToken == "" {
		return "", errors.New("invalid split tag")
	}

	return splitToken, nil
}

func (p *parser) parseDefaultOption() (string, error) {
	// consume the default token
	p.currentIndex++

	// now we should see an = sign
	if ct := p.currentToken(); ct != "=" {
		return "", errors.New("invalid split tag")
	}
	p.currentIndex++

	// the token now will be our start token
	start := p.currentToken()
	switch start {
	case "", ",":
		return "", errors.New("invalid default tag")

	case "[":
		// if we see a [ we need to consume until we see a ]
		startingIndex := p.currentIndex
		for {
			// if we run out of token before seeing a ], it is invalid
			if p.currentIndex >= len(p.tokens) {
				return "", errors.New("invalid default tag")
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
			return MinienvTag{}, fmt.Errorf("invalid token \"%s\"", token)
		}

		p.currentIndex++ // move to the next token
	}

	return tag, nil
}
