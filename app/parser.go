package main

import "strings"

func parse(input string) []string {
	var tokens []string
	var current strings.Builder
	var inSingleQuotes, inDoubleQuotes, escaped bool
	for _, r := range input {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
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
