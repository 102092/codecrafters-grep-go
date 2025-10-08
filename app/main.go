package main

import (
	"fmt"
	"io"
	"os"
	"strings"
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
	// Special case: pattern starts with ^ and ends with $, must match the whole line
	if strings.HasPrefix(pattern, "^") && strings.HasSuffix(pattern, "$") {
		pattern = pattern[1 : len(pattern)-1] // Remove both ^ and $, ^apple$ -> apple
		return string(inputText) == pattern, nil
	}

	// Special case: pattern starts with ^, must match from the beginning
	if strings.HasPrefix(pattern, "^") {
		pattern = pattern[1:] // Remove leading ^, ^apple -> apple
		return strings.HasPrefix(string(inputText), pattern), nil
	}

	// Special case: pattern ends with $, must match at the end
	if strings.HasSuffix(pattern, "$") {
		pattern = pattern[0 : len(pattern)-1] // Remove trailing $, apple$ -> apple
		return strings.HasSuffix(string(inputText), pattern), nil
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
// It handles quantifiers like + (one or more) using backtracking.
// It returns true if all tokens match consecutively from the start position.
func matchFromPosition(inputText []byte, tokens []Token, startIndex int) bool {
	return matchFromPositionRecursive(inputText, tokens, 0, startIndex)
}

// matchFromPositionRecursive recursively matches tokens with backtracking support.
// It tries different match lengths for quantified tokens and backtracks on failure.
//
// Parameters:
//   - inputText: the input string to match against
//   - tokens: the parsed pattern tokens
//   - tokenIndex: current token being matched
//   - inputIndex: current position in the input text
//
// Returns true if all remaining tokens can be matched from the current position.
func matchFromPositionRecursive(inputText []byte, tokens []Token, tokenIndex int, inputIndex int) bool {
	// Base case: all tokens matched successfully
	if tokenIndex >= len(tokens) {
		return true
	}

	token := tokens[tokenIndex]

	if token.Quantifier == OneOrMore {
		// + quantifier: match one or more times with backtracking
		// Must match at least once
		if inputIndex >= len(inputText) || !matchToken(token, inputText[inputIndex]) {
			return false
		}

		// Try matching 1, 2, 3, ... times (backtracking)
		// Start with minimum (1) and try increasing matches
		for i := inputIndex; i < len(inputText) && matchToken(token, inputText[i]); i++ {
			// Try matching the rest of the pattern with current match count
			if matchFromPositionRecursive(inputText, tokens, tokenIndex+1, i+1) {
				return true // Found a successful match
			}
			// If failed, backtrack and try matching one more character
		}

		return false // All attempts failed

	} else {
		// No quantifier: match exactly once
		if inputIndex >= len(inputText) {
			return false
		}

		if matchToken(token, inputText[inputIndex]) {
			// Recursively match the rest of the pattern
			return matchFromPositionRecursive(inputText, tokens, tokenIndex+1, inputIndex+1)
		}

		return false // Token doesn't match
	}
}
