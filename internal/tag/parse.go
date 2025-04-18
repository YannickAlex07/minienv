package tag

import (
	"errors"
	"fmt"
	"iter"
)

// TOKENIZATION

func tokenize(tag string) ([]string, error) {
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

	return tokens, nil
}

// PARSER

type tagParser struct {
	tokens       []string
	currentIndex int
}

func newTagParser(tag string) (*tagParser, error) {
	tokens, err := tokenize(tag)
	if err != nil {
		return nil, err
	}

	return &tagParser{
		tokens:       tokens,
		currentIndex: 0,
	}, nil
}

func (p *tagParser) Parse() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for p.currentIndex < len(p.tokens) {
			currentToken := p.tokens[p.currentIndex]

			switch currentToken {
			case "split":
				s, err := p.parseSplit()
				if !yield(s, err) {
					return
				}

			case "default":
				s, err := p.parseDefault()
				if !yield(s, err) {
					return
				}

			case "optional":
				p.currentIndex++
				if !yield(currentToken, nil) {
					return
				}
			case ",":
				p.currentIndex++
			default:
				// if we don't recognize the token, something is invalid
				if !yield("", errors.New("invalid token")) {
					return
				}
			}
		}
	}
}

func (p *tagParser) parseSplit() (string, error) {
	// cosume the split token
	p.currentIndex++

	// check and consume the =
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

	return fmt.Sprintf("split=%s", splitOn), nil
}

func (p *tagParser) parseDefault() (string, error) {
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
	// if the current token is a [, we need to consume until we find a ]
	if p.tokens[p.currentIndex] == "[" {
		// consume the opening bracket
		value += p.tokens[p.currentIndex]
		p.currentIndex++

		for p.tokens[p.currentIndex] != "]" {
			value += p.tokens[p.currentIndex]
			p.currentIndex++

			if p.isEOF() {
				return "", errors.New("invalid default tag")
			}
		}

		// consume the closing bracket
		value += p.tokens[p.currentIndex]
		p.currentIndex++
	} else {
		value = p.tokens[p.currentIndex]
		p.currentIndex++
	}

	return fmt.Sprintf("default=%s", value), nil
}

func (p *tagParser) isEOF() bool {
	return p.currentIndex >= len(p.tokens)
}

func (p *tagParser) peek(n int) (string, bool) {
	if p.currentIndex+n >= len(p.tokens) {
		return "", false
	}

	return p.tokens[p.currentIndex+n], true
}
