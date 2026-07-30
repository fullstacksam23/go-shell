package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

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
		insertTrie(k)
	}
	loadExecutables()
}

func insertTrie(elem string) {
	node := currTrie.root
	for _, r := range elem {
		if node.children[r] == nil {
			node.children[r] = &TrieNode{
				children: make(map[rune]*TrieNode),
			}
		}
		node = node.children[r]
	}
	node.isEnd = true
}

func loadExecutables() {
	pathEnv := os.Getenv("PATH")

	for _, dir := range filepath.SplitList(pathEnv) {
		dirs, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range dirs {
			if entry.IsDir() {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Ignore non-executables
			if runtime.GOOS == "windows" {
				if filepath.Ext(entry.Name()) != ".exe" {
					continue
				}
			} else {
				if info.Mode()&0111 == 0 {
					continue
				}
			}

			insertTrie(entry.Name())
		}
	}
}

func autoCompleteCommand(line string) ([]rune, []string) {
	curr := currTrie.root

	for _, r := range line {
		next, ok := curr.children[r]
		if !ok {
			return nil, nil
		}
		curr = next
	}
	node := curr
	matches := dfs(curr, []rune{}, nil)
	slices.Sort(matches)

	longestCommon := make([]rune, 0)
	for len(node.children) == 1 && !node.isEnd {
		var nxt *TrieNode
		for k, v := range node.children {
			longestCommon = append(longestCommon, k)
			nxt = v
		}
		node = nxt
	}
	return longestCommon, matches
}

func dfs(node *TrieNode, curr []rune, matches []string) []string {
	if node.isEnd {
		matches = append(matches, string(curr))
	}

	for r, child := range node.children {
		curr = append(curr, r)
		matches = dfs(child, curr, matches)
		curr = curr[:len(curr)-1]
	}

	return matches
}

func autoCompleteFileName(path, word string) []string {
	matches := []string{}
	dirs, err := os.ReadDir(path)
	if err != nil {
		fmt.Fprintln(crlfWriter{os.Stderr}, err)
		return matches
	}
	for _, d := range dirs {

		if strings.HasPrefix(d.Name(), word) {
			if d.IsDir() {
				matches = append(matches, d.Name()[len(word):]+"/")
			} else {
				matches = append(matches, d.Name()[len(word):]+" ")
			}
		}

	}
	return matches
}
