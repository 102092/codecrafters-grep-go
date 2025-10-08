package main

import (
	"testing"
)

func TestParseTokens(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []Token
		wantErr bool
	}{
		// Literal characters
		{
			name:    "single literal character",
			pattern: "a",
			want: []Token{
				{Type: Literal, Value: "a", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "multiple literal characters",
			pattern: "abc",
			want: []Token{
				{Type: Literal, Value: "a", Quantifier: None},
				{Type: Literal, Value: "b", Quantifier: None},
				{Type: Literal, Value: "c", Quantifier: None},
			},
			wantErr: false,
		},
		// Escape sequences
		{
			name:    "\\d digit pattern",
			pattern: "\\d",
			want: []Token{
				{Type: Digit, Value: "\\d", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "\\w word pattern",
			pattern: "\\w",
			want: []Token{
				{Type: Word, Value: "\\w", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "\\\\ literal backslash",
			pattern: "\\\\",
			want: []Token{
				{Type: Literal, Value: "\\", Quantifier: None},
			},
			wantErr: false,
		},
		// Character classes
		{
			name:    "[abc] positive character class",
			pattern: "[abc]",
			want: []Token{
				{Type: CharClass, Value: "abc", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "[^abc] negative character class",
			pattern: "[^abc]",
			want: []Token{
				{Type: NegCharClass, Value: "abc", Quantifier: None},
			},
			wantErr: false,
		},
		// Quantifiers
		{
			name:    "a+ with quantifier",
			pattern: "a+",
			want: []Token{
				{Type: Literal, Value: "a", Quantifier: OneOrMore},
			},
			wantErr: false,
		},
		{
			name:    "\\d+ with quantifier",
			pattern: "\\d+",
			want: []Token{
				{Type: Digit, Value: "\\d", Quantifier: OneOrMore},
			},
			wantErr: false,
		},
		{
			name:    "[abc]+ character class with quantifier",
			pattern: "[abc]+",
			want: []Token{
				{Type: CharClass, Value: "abc", Quantifier: OneOrMore},
			},
			wantErr: false,
		},
		{
			name:    "a? with zero-or-one quantifier",
			pattern: "a?",
			want: []Token{
				{Type: Literal, Value: "a", Quantifier: ZeroOrOne},
			},
			wantErr: false,
		},
		{
			name:    "\\d? with zero-or-one quantifier",
			pattern: "\\d?",
			want: []Token{
				{Type: Digit, Value: "\\d", Quantifier: ZeroOrOne},
			},
			wantErr: false,
		},
		{
			name:    "[abc]? character class with zero-or-one quantifier",
			pattern: "[abc]?",
			want: []Token{
				{Type: CharClass, Value: "abc", Quantifier: ZeroOrOne},
			},
			wantErr: false,
		},
		// Complex patterns
		{
			name:    "ca+ts pattern",
			pattern: "ca+ts",
			want: []Token{
				{Type: Literal, Value: "c", Quantifier: None},
				{Type: Literal, Value: "a", Quantifier: OneOrMore},
				{Type: Literal, Value: "t", Quantifier: None},
				{Type: Literal, Value: "s", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "\\d+ apple pattern",
			pattern: "\\d+ apple",
			want: []Token{
				{Type: Digit, Value: "\\d", Quantifier: OneOrMore},
				{Type: Literal, Value: " ", Quantifier: None},
				{Type: Literal, Value: "a", Quantifier: None},
				{Type: Literal, Value: "p", Quantifier: None},
				{Type: Literal, Value: "p", Quantifier: None},
				{Type: Literal, Value: "l", Quantifier: None},
				{Type: Literal, Value: "e", Quantifier: None},
			},
			wantErr: false,
		},
		{
			name:    "a\\\\b pattern (literal backslash between chars)",
			pattern: "a\\\\b",
			want: []Token{
				{Type: Literal, Value: "a", Quantifier: None},
				{Type: Literal, Value: "\\", Quantifier: None},
				{Type: Literal, Value: "b", Quantifier: None},
			},
			wantErr: false,
		},
		// Error cases
		{
			name:    "unclosed character class",
			pattern: "[abc",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "+ without preceding character",
			pattern: "+abc",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unsupported escape sequence",
			pattern: "\\x",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTokens(tt.pattern)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalTokenSlices(got, tt.want) {
				t.Errorf("parseTokens() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestMatchToken(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		b     byte
		want  bool
	}{
		// Literal tokens
		{
			name:  "literal 'a' matches 'a'",
			token: Token{Type: Literal, Value: "a"},
			b:     'a',
			want:  true,
		},
		{
			name:  "literal 'a' does not match 'b'",
			token: Token{Type: Literal, Value: "a"},
			b:     'b',
			want:  false,
		},
		{
			name:  "literal '\\' matches '\\'",
			token: Token{Type: Literal, Value: "\\"},
			b:     '\\',
			want:  true,
		},
		{
			name:  "literal '\\' does not match 'a'",
			token: Token{Type: Literal, Value: "\\"},
			b:     'a',
			want:  false,
		},
		// Digit tokens
		{
			name:  "\\d matches '5'",
			token: Token{Type: Digit, Value: "\\d"},
			b:     '5',
			want:  true,
		},
		{
			name:  "\\d does not match 'a'",
			token: Token{Type: Digit, Value: "\\d"},
			b:     'a',
			want:  false,
		},
		// Word tokens
		{
			name:  "\\w matches 'a'",
			token: Token{Type: Word, Value: "\\w"},
			b:     'a',
			want:  true,
		},
		{
			name:  "\\w matches '5'",
			token: Token{Type: Word, Value: "\\w"},
			b:     '5',
			want:  true,
		},
		{
			name:  "\\w matches '_'",
			token: Token{Type: Word, Value: "\\w"},
			b:     '_',
			want:  true,
		},
		{
			name:  "\\w does not match '!'",
			token: Token{Type: Word, Value: "\\w"},
			b:     '!',
			want:  false,
		},
		// CharClass tokens
		{
			name:  "[abc] matches 'a'",
			token: Token{Type: CharClass, Value: "abc"},
			b:     'a',
			want:  true,
		},
		{
			name:  "[abc] matches 'c'",
			token: Token{Type: CharClass, Value: "abc"},
			b:     'c',
			want:  true,
		},
		{
			name:  "[abc] does not match 'z'",
			token: Token{Type: CharClass, Value: "abc"},
			b:     'z',
			want:  false,
		},
		// NegCharClass tokens
		{
			name:  "[^abc] matches 'z'",
			token: Token{Type: NegCharClass, Value: "abc"},
			b:     'z',
			want:  true,
		},
		{
			name:  "[^abc] does not match 'a'",
			token: Token{Type: NegCharClass, Value: "abc"},
			b:     'a',
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchToken(tt.token, tt.b)
			if got != tt.want {
				t.Errorf("matchToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two slices of tokens
func equalTokenSlices(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Type != b[i].Type || a[i].Value != b[i].Value || a[i].Quantifier != b[i].Quantifier {
			return false
		}
	}
	return true
}
