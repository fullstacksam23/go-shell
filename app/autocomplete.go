package main

import (
	"fmt"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/builtins"
)

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
}

type Trie struct {
	root *TrieNode
}

var currTrie *Trie

func SetupAutoComplete() {
	currTrie = &Trie{
		root: &TrieNode{
			children: make(map[rune]*TrieNode),
		},
	}

	for k := range builtins.BuiltIn {
		node := currTrie.root
		for _, r := range k {
			if node.children[r] == nil {
				node.children[r] = &TrieNode{
					children: make(map[rune]*TrieNode),
				}
			}
			node = node.children[r]
		}
		node.isEnd = true
	}
}

func autocompleteSuffix(line string) string {
	curr := currTrie.root

	for _, r := range line {
		if curr.children[r] == nil {
			return ""
		}
		curr = curr.children[r]
	}

	if curr.isEnd {
		return " "
	}

	return dfs(curr) + " "
}

func dfs(node *TrieNode) string {
	if node.isEnd {
		return ""
	}

	for r, child := range node.children {
		return string(r) + dfs(child)
	}

	return ""
}

type TrieCompleter struct{}

func (t *TrieCompleter) Do(line []rune, pos int) ([][]rune, int) {
	prefix := string(line[:pos])

	suffix := autocompleteSuffix(prefix)
	if suffix == "" {
		fmt.Print("\x07")
		return nil, 0
	}

	return [][]rune{[]rune(suffix)}, 0
}

var completer readline.AutoCompleter = &TrieCompleter{}
