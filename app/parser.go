package main

import (
	"errors"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/core"
)

func parse(tokens []string) (*core.ParsedCommand, error) {
	if len(tokens) < 1 {
		return nil, errors.New("missing input")
	}
	p := core.ParsedCommand{
		Command: tokens[0],
	}
	redirect := -1
	for i, v := range tokens {
		if v == ">" || v == "1>" {
			redirect = i
			break
		}
	}
	if redirect == -1 {
		p.Args = tokens[1:]
		return &p, nil
	}
	if redirect+1 >= len(tokens) {
		return nil, errors.New("echo: missing file operand")
	}

	p.Args = tokens[1:redirect]
	p.StdoutFile = tokens[redirect+1]
	p.StdoutRedirect = true
	return &p, nil

}

func tokenize(input string) []string {
	var tokens []string
	var current strings.Builder
	var inSingleQuotes, inDoubleQuotes, escaped bool
	for _, r := range input {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case !inSingleQuotes && r == '\\':
			escaped = true

		case r == '"' && !inDoubleQuotes && !inSingleQuotes:
			inDoubleQuotes = true

		case r == '"' && inDoubleQuotes && !inSingleQuotes:
			inDoubleQuotes = false

		case r == '\'' && !inSingleQuotes && !inDoubleQuotes:
			inSingleQuotes = true

		case r == '\'' && inSingleQuotes && !inDoubleQuotes:
			inSingleQuotes = false

		case r == ' ' && !inSingleQuotes && !inDoubleQuotes:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}
