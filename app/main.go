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
// It handles quantifiers like + (one or more) using greedy matching.
// It returns true if all tokens match consecutively from the start position.
func matchFromPosition(inputText []byte, tokens []Token, startIndex int) bool {
	inputIndex := startIndex

	// Try to match each token in sequence
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		// Handle quantifiers
		if token.Quantifier == OneOrMore {
			// + quantifier: match one or more times (greedy)
			count := 0

			// Consume as many matching characters as possible
			for inputIndex < len(inputText) && matchToken(token, inputText[inputIndex]) {
				count++
				inputIndex++
			}

			// Must match at least once
			if count < 1 {
				return false
			}

		} else {
			// No quantifier: match exactly once
			if inputIndex >= len(inputText) {
				return false
			}

			if matchToken(token, inputText[inputIndex]) {
				inputIndex++ // Move to next character
			} else {
				return false // Token doesn't match, fail immediately
			}
		}
	}

	return true // All tokens matched successfully
}
