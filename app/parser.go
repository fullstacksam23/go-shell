package main

import "regexp"

var tokenRE = regexp.MustCompile(`"[^"]*"|\S+`)

func parse(input string) []string {
	tokens := tokenRE.FindAllString(input, -1)

	for i, t := range tokens {
		if len(t) >= 2 && t[0] == '"' && t[len(t)-1] == '"' {
			tokens[i] = t[1 : len(t)-1]
		}
	}

	return tokens
}
