package main

import (
	"fmt"
	"strings"
)

const (
	digits    = "0123456789"
	wordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// TokenType represents the type of a pattern token.
type TokenType int

const (
	Literal      TokenType = iota // Single literal character: "a", "b"
	Digit                         // \d - digit character
	Word                          // \w - word character
	CharClass                     // [abc] - positive character class
	NegCharClass                  // [^abc] - negative character class
)

// QuantifierType represents the quantifier applied to a token.
type QuantifierType int

const (
	None      QuantifierType = iota // Exactly one (no quantifier)
	OneOrMore                       // + (one or more)
)

// Token represents a single pattern matching unit with optional quantifier.
type Token struct {
	Type       TokenType      // Type of the token
	Value      string         // The pattern value (e.g., "a", "\\d", "abc" for char class)
	Quantifier QuantifierType // Quantifier type
}

// parseTokens breaks a pattern string into individual tokens with quantifiers.
//
// Supported tokens:
//   - Escape sequences: \d, \w
//   - Character classes: [abc], [^abc]
//   - Single characters: any other character
//   - Quantifiers: + (one or more)
//
// It returns an error if:
//   - A character class is not properly closed with ']'
//   - A quantifier appears without a preceding character
func parseTokens(pattern string) ([]Token, error) {
	var tokens []Token

	for i := 0; i < len(pattern); {
		var token Token
		advance := 1

		// Escape sequences: \d, \w, etc.
		if pattern[i] == '\\' && i+1 < len(pattern) {
			advance = 2

			// Create token based on escape sequence type
			switch pattern[i+1] {
			case 'd':
				token = Token{
					Type:       Digit,
					Value:      "\\d",
					Quantifier: None,
				}
			case 'w':
				token = Token{
					Type:       Word,
					Value:      "\\w",
					Quantifier: None,
				}
			case '\\':
				// Literal backslash: \\ represents a single '\'
				token = Token{
					Type:       Literal,
					Value:      "\\",
					Quantifier: None,
				}
			default:
				return nil, fmt.Errorf("unsupported escape sequence: %s", pattern[i:i+2])
			}

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

			advance = j + 1 - i
			patternValue := pattern[i : j+1]

			// Check if it's a negative character class
			if len(patternValue) >= 3 && patternValue[1] == '^' {
				token = Token{
					Type:       NegCharClass,
					Value:      patternValue[2 : len(patternValue)-1], // Extract "abc" from "[^abc]"
					Quantifier: None,
				}
			} else {
				token = Token{
					Type:       CharClass,
					Value:      patternValue[1 : len(patternValue)-1], // Extract "abc" from "[abc]"
					Quantifier: None,
				}
			}

			// Single literal character
		} else {
			// Check for invalid + at the beginning or after special chars
			if pattern[i] == '+' {
				return nil, fmt.Errorf("invalid pattern: + must follow a character at position %d", i)
			}

			token = Token{
				Type:       Literal,
				Value:      string(pattern[i]),
				Quantifier: None,
			}
		}

		// Check for quantifier after the current token
		advance += parseQuantifierIfPresent(pattern, i+advance, &token)

		tokens = append(tokens, token)
		i += advance
	}

	return tokens, nil
}

// parseQuantifierIfPresent checks for a quantifier (+) after the current position
// and updates the token accordingly. Returns number of characters consumed.
//
// Parameters:
//   - pattern: the pattern string being parsed
//   - pos: current position in the pattern to check for quantifier
//   - token: pointer to the token to update with quantifier
//
// Returns the number of characters to advance (0 if no quantifier, 1 if quantifier found).
func parseQuantifierIfPresent(pattern string, pos int, token *Token) int {
	if pos < len(pattern) && pattern[pos] == '+' {
		token.Quantifier = OneOrMore
		return 1
	}
	return 0
}

// matchToken checks if a Token matches a single byte.
// It returns true if the byte matches the token's pattern.
func matchToken(token Token, b byte) bool {
	switch token.Type {
	case Literal:
		return token.Value == string(b)

	case Digit:
		return strings.ContainsAny(string(b), digits)

	case Word:
		return strings.ContainsAny(string(b), wordChars)

	case CharClass:
		if token.Value == "" {

			return false
		}
		return strings.ContainsRune(token.Value, rune(b))

	case NegCharClass:
		if token.Value == "" {
			return false
		}
		return !strings.ContainsRune(token.Value, rune(b))

	default:
		return false
	}
}
