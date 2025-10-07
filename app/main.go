package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	digits    = "0123456789"
	wordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

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

// matchLine checks if the pattern matches anywhere in the line.
// It tries matching from every position in the line until a match is found.
func matchLine(inputText []byte, pattern string) (bool, error) {
	// Special case: pattern starts with ^, must match from the beginning
	if strings.HasPrefix(pattern, "^") {
		pattern = pattern[1:] // Remove leading ^, ^apple -> apple
		return strings.HasPrefix(string(inputText), pattern), nil
	}

	tokens, err := parseTokens(pattern)
	if err != nil {
		return false, err
	}

	// Try matching from each position in the inputText
	for start := 0; start < len(inputText); start++ {
		if matchFromPosition(inputText, tokens, start) {
			return true, nil
		}
	}

	return false, nil
}

// matchFromPosition attempts to match all tokens sequentially starting from the given position.
// It returns true if all tokens match consecutively from the start position.
func matchFromPosition(inputText []byte, tokens []string, startIndex int) bool {
	inputIndex := startIndex

	// Try to match each token in sequence
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Check if we've run out of inputIndex characters
		if inputIndex >= len(inputText) {
			return false
		}

		if matchSingleToken(token, inputText[inputIndex]) {
			inputIndex++ // Move to next character
		} else {
			return false // Token doesn't match, fail immediately
		}
	}

	return true // All tokens matched successfully
}

// matchSingleToken checks if a single token matches a single byte.
//
// Supported token types:
//   - \d: digit characters (0-9)
//   - \w: word characters (letters, digits, underscore)
//   - [abc]: positive character groups (matches if b is in the set)
//   - [^abc]: negative character groups (matches if b is NOT in the set)
//   - literal characters: exact match
func matchSingleToken(token string, b byte) bool {
	// \d: digit character (0-9)
	if token == "\\d" {
		return strings.ContainsAny(string(b), digits)
	}

	// \w: word character (letters, digits, underscore)
	if token == "\\w" {
		return strings.ContainsAny(string(b), wordChars)
	}

	// [^abc]: negative character group - matches if b is NOT in the set
	if strings.HasPrefix(token, "[^") && strings.HasSuffix(token, "]") {
		charClass := token[2 : len(token)-1]
		if charClass == "" {
			return false
		}

		return !strings.ContainsRune(charClass, rune(b))
	}

	// [abc]: positive character group - matches if b is in the set
	if strings.HasPrefix(token, "[") && strings.HasSuffix(token, "]") {
		charClass := token[1 : len(token)-1]
		if charClass == "" {
			return false
		}

		return strings.ContainsRune(charClass, rune(b))
	}

	// Literal character: exact match
	return token == string(b)
}

// parseTokens breaks a pattern string into individual tokens.
// Tokens can be:
//   - Escape sequences: \d, \w (2 characters treated as one token)
//   - Character classes: [abc], [^abc] (entire bracket expression as one token)
//   - Single characters: any other character
//
// Returns an error if a character class is not properly closed with ']'.
func parseTokens(pattern string) ([]string, error) {
	var tokens []string

	for i := 0; i < len(pattern); {
		// Escape sequences: \d, \w, etc.
		if pattern[i] == '\\' && i+1 < len(pattern) {
			tokens = append(tokens, pattern[i:i+2])
			i += 2

			// Character classes: [abc] or [^abc]
		} else if pattern[i] == '[' {
			j := i + 1
			// Find the closing ']'
			for j < len(pattern) && pattern[j] != ']' {
				j++
			}
			// Check if we found a closing bracket
			if j >= len(pattern) {
				return nil, fmt.Errorf("unclosed character class starting at position %d", i)
			}

			tokens = append(tokens, pattern[i:j+1])
			i = j + 1

			// Single literal character
		} else {
			tokens = append(tokens, string(pattern[i]))
			i++
		}
	}
	return tokens, nil
}
