package tag

import (
	"errors"
	"fmt"
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

func (p *parser) Parse() (MinienvTag, error) {
	tag := MinienvTag{}

	for p.currentIndex < len(p.tokens) {
		currentToken := p.tokens[p.currentIndex]

		// the first token is always the lookup name
		if p.currentIndex == 0 {
			tag.LookupName = currentToken
			p.currentIndex++
			continue
		}

		switch currentToken {
		case "optional":
			tag.Optional = true
			p.currentIndex++

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
			p.currentIndex++

		default:
			// if we don't recognize the token, something is invalid
			return MinienvTag{}, fmt.Errorf("invalid token \"%s\"", currentToken)
		}
	}

	return tag, nil
}

func (p *parser) parseSplitOption() (string, error) {
	// cosume the split token itself
	p.currentIndex++

	// check and consume the = token
	if p.isEOF() || p.tokens[p.currentIndex] != "=" {
		return "", errors.New("invalid split tag")
	}
	p.currentIndex++

	if p.isEOF() {
		return "", errors.New("invalid split tag")
	}

	// if after the current token there is still more and it isn't a ",", the tag is invalid
	peeked, ok := p.peek(1)
	if ok && peeked != "," {
		return "", errors.New("invalid split tag")
	}

	// our current token will be the split character and we can consume it
	splitOn := p.tokens[p.currentIndex]
	p.currentIndex++

	return splitOn, nil
}

func (p *parser) parseDefaultOption() (string, error) {
	// consume the default token
	p.currentIndex++

	// check and consume the =
	if p.isEOF() || p.tokens[p.currentIndex] != "=" {
		return "", errors.New("invalid default tag")
	}
	p.currentIndex++

	// if we are EOF, this is invalid
	if p.isEOF() {
		return "", errors.New("invalid default tag")
	}

	value := ""
	switch p.tokens[p.currentIndex] {
	case ",":
		// a comma is not a valid default value currently
		return "", errors.New("invalid default tag")

	case "[":
		// if the current token is a [, we need to consume until we find a ]
		// consume the opening bracket
		p.currentIndex++

		for p.tokens[p.currentIndex] != "]" {
			value += p.tokens[p.currentIndex]
			p.currentIndex++

			if p.isEOF() {
				return "", errors.New("invalid default tag")
			}
		}

		// consume the closing bracket
		p.currentIndex++

	default:
		value = p.tokens[p.currentIndex]
		p.currentIndex++
	}

	return value, nil
}

func (p *parser) isEOF() bool {
	return p.currentIndex >= len(p.tokens)
}

func (p *parser) peek(n int) (string, bool) {
	if p.currentIndex+n >= len(p.tokens) {
		return "", false
	}

	return p.tokens[p.currentIndex+n], true
}
