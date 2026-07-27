package main

import "github.com/codecrafters-io/shell-starter-go/builtins"

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

	for k, _ := range builtins.BuiltIn {
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

func autocomplete(line string) string {
	curr := currTrie.root

	for i, r := range line {
		// no match
		if curr.children[r] == nil {
			return line[:i+1]
		}
		curr = curr.children[r]
	}
	//check if the user typed the whole command
	if curr.isEnd {
		return line + " "
	}
	return line + dfs(curr) + " "
}

func dfs(node *TrieNode) string {
	if node.isEnd {
		return ""
	}
	for r, v := range node.children {
		return string(r) + dfs(v)
	}
	return ""
}
