package main

import "strings"

func parse(input string) []string {
	var tokens []string
	var current strings.Builder
	var inSingleQuotes bool
	for _, r := range input {
		switch {
		case r == '\'' && !inSingleQuotes:
			inSingleQuotes = true

		case r == '\'' && inSingleQuotes:
			inSingleQuotes = false

		case r == ' ' && !inSingleQuotes:
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
