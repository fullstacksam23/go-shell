package main

import (
	"errors"
	"strings"
)

type ParsedCommand struct {
	Command        string
	Args           []string
	StdoutFile     string
	StdoutRedirect bool
	StdoutAppend   bool
	StderrFile     string
	StdErrRedirect bool
	StdErrAppend   bool
}

func parse(tokens []string) (*ParsedCommand, error) {
	if len(tokens) < 1 {
		return nil, errors.New("missing input")
	}
	p := ParsedCommand{
		Command: tokens[0],
	}
	stdoutRedirect := -1
	stderrRedirect := -1
	for i, v := range tokens {
		if v == ">" || v == "1>" || v == ">>" || v == "1>>" {
			stdoutRedirect = i
		}
		if v == "2>" || v == "2>>" {
			stderrRedirect = i
		}
	}

	if stdoutRedirect == -1 {
		p.Args = tokens[1:]
	} else {

		if stdoutRedirect+1 >= len(tokens) {
			return nil, errors.New("missing file operand")
		}

		p.Args = tokens[1:stdoutRedirect]
		p.StdoutFile = tokens[stdoutRedirect+1]
		if tokens[stdoutRedirect] == ">>" || tokens[stdoutRedirect] == "1>>" {
			p.StdoutAppend = true
		} else {
			p.StdoutRedirect = true
		}
	}

	if stderrRedirect != -1 {
		if stdoutRedirect == -1 || stderrRedirect < stdoutRedirect {
			p.Args = tokens[1:stderrRedirect]
		}
		if stderrRedirect+1 >= len(tokens) {
			return nil, errors.New("missing file operand")
		}
		p.StderrFile = tokens[stderrRedirect+1]

		if tokens[stderrRedirect] == "2>" {
			p.StdErrRedirect = true
		} else {
			p.StdErrAppend = true
		}
	}
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
