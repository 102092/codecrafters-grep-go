package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "bytes" import above (feel free to remove this!)
var _ = bytes.ContainsAny

// Usage: echo <input_text> | your_program.sh -E <pattern>
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2) // 1 means no lines were selected, >1 means error
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}

	// default exit code is 0 which means success
}

func matchLine(line []byte, pattern string) (bool, error) {
	// Note: "\\d" represents the literal string "\d" (backslash needs to be escaped in Go string literals)
	tokens := parseTokens(pattern)

	for start := 0; start < len(line); start++ {
		if matchFromPosition(line, tokens, start) {
			return true, nil
		}
	}

	return false, nil
}

func matchFromPosition(line []byte, tokens []string, start int) bool {
	input := start

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		if input >= len(line) {
			return false
		}

		if matchSingleToken(token, line[input]) {
			input++
		} else {
			return false
		}
	}

	return true
}

func matchSingleToken(token string, b byte) bool {
	// digit character: 0-9
	if token == "\\d" {
		return strings.ContainsAny(string(b), "0123456789")
	}

	// word character: letters, digits, underscore
	if token == "\\w" {
		return strings.ContainsAny(string(b), "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	}

	// Negated character group: [^abc]
	if strings.HasPrefix(token, "[^") && strings.HasSuffix(token, "]") {
		char := token[2 : len(token)-1]
		if char == "" {
			return false
		}

		return !strings.ContainsRune(char, rune(b))
	}

	// Positive character group: [abc]
	if strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]") {
		char := token[1 : len(token)-1]
		if char == "" {
			return false
		}

		return strings.ContainsRune(char, rune(b))
	}

	return token == string(b)
}

func parseTokens(pattern string) []string {
	tokens := []string{}

	for i := 0; i < len(pattern); {
		if pattern[i] == '\\' && i+1 < len(pattern) {
			tokens = append(tokens, pattern[i:i+2])
			i += 2

		} else if pattern[i] == '[' {
			j := i + 1
			for j < len(pattern) && pattern[j] != ']' {
				j++
			}
			tokens = append(tokens, pattern[i:j+1])
			i = j + 1
		} else {
			tokens = append(tokens, string(pattern[i]))
			i++
		}
	}
	return tokens
}
