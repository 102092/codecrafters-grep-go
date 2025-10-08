package main

import (
	"testing"
)

func TestMatchLine(t *testing.T) {
	tests := []struct {
		name    string
		line    []byte
		pattern string
		want    bool
		wantErr bool
	}{
		// \d 패턴 테스트
		{
			name:    "\\d matches digit",
			line:    []byte("hello123"),
			pattern: "\\d",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d no match",
			line:    []byte("hello"),
			pattern: "\\d",
			want:    false,
			wantErr: false,
		},
		// \w 패턴 테스트
		{
			name:    "\\w matches word char",
			line:    []byte("hello"),
			pattern: "\\w",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\w no match",
			line:    []byte("!!!"),
			pattern: "\\w",
			want:    false,
			wantErr: false,
		},
		// 단일 문자 테스트
		{
			name:    "single char match",
			line:    []byte("apple"),
			pattern: "a",
			want:    true,
			wantErr: false,
		},
		{
			name:    "single char no match",
			line:    []byte("apple"),
			pattern: "z",
			want:    false,
			wantErr: false,
		},
		// Positive character groups 테스트
		{
			name:    "[abc] matches a",
			line:    []byte("apple"),
			pattern: "[abc]",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc] matches b",
			line:    []byte("banana"),
			pattern: "[abc]",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc] no match",
			line:    []byte("xyz"),
			pattern: "[abc]",
			want:    false,
			wantErr: false,
		},
		{
			name:    "[xyz] matches z",
			line:    []byte("buzz"),
			pattern: "[xyz]",
			want:    true,
			wantErr: false,
		},
		// Negative character groups 테스트
		{
			name:    "[^abc] matches cat (has t)",
			line:    []byte("cat"),
			pattern: "[^abc]",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[^abc] does not match cab",
			line:    []byte("cab"),
			pattern: "[^abc]",
			want:    false,
			wantErr: false,
		},
		{
			name:    "[^ab] matches abc (has c)",
			line:    []byte("abc"),
			pattern: "[^ab]",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[^xyz] matches hello",
			line:    []byte("hello"),
			pattern: "[^xyz]",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[^xyz] does not match xyz",
			line:    []byte("xyz"),
			pattern: "[^xyz]",
			want:    false,
			wantErr: false,
		},
		// Pattern sequence tests
		{
			name:    "\\d apple matches 1 apple",
			line:    []byte("1 apple"),
			pattern: "\\d apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d apple does not match 1 orange",
			line:    []byte("1 orange"),
			pattern: "\\d apple",
			want:    false,
			wantErr: false,
		},
		{
			name:    "\\d\\d\\d apple matches 100 apples",
			line:    []byte("100 apples"),
			pattern: "\\d\\d\\d apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d\\d\\d apple does not match 1 apple",
			line:    []byte("1 apple"),
			pattern: "\\d\\d\\d apple",
			want:    false,
			wantErr: false,
		},
		{
			name:    "\\d \\w\\w\\ws matches 3 dogs",
			line:    []byte("3 dogs"),
			pattern: "\\d \\w\\w\\ws",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d \\w\\w\\ws matches 4 cats",
			line:    []byte("4 cats"),
			pattern: "\\d \\w\\w\\ws",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d \\w\\w\\ws does not match 1 dog",
			line:    []byte("1 dog"),
			pattern: "\\d \\w\\w\\ws",
			want:    false,
			wantErr: false,
		},
		{
			name:    "multi-char sequence abc",
			line:    []byte("xyzabcdef"),
			pattern: "abc",
			want:    true,
			wantErr: false,
		},
		// Start of string anchor tests
		{
			name:    "^apple matches apple pie",
			line:    []byte("apple pie"),
			pattern: "^apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "^apple does not match green apple",
			line:    []byte("green apple"),
			pattern: "^apple",
			want:    false,
			wantErr: false,
		},
		{
			name:    "^apple matches apple exactly",
			line:    []byte("apple"),
			pattern: "^apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "^log matches log message",
			line:    []byte("log message"),
			pattern: "^log",
			want:    true,
			wantErr: false,
		},
		// End of string anchor tests
		{
			name:    "apple$ matches green apple",
			line:    []byte("green apple"),
			pattern: "apple$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "apple$ does not match apple pie",
			line:    []byte("apple pie"),
			pattern: "apple$",
			want:    false,
			wantErr: false,
		},
		{
			name:    "apple$ matches apple exactly",
			line:    []byte("apple"),
			pattern: "apple$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "world$ matches hello world",
			line:    []byte("hello world"),
			pattern: "world$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "strawberry$ matches orange_strawberry",
			line:    []byte("orange_strawberry"),
			pattern: "strawberry$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "strawberry$ does not match strawberry_orange",
			line:    []byte("strawberry_orange"),
			pattern: "strawberry$",
			want:    false,
			wantErr: false,
		},
		// Both anchors: ^...$ (exact match)
		{
			name:    "^apple$ matches apple exactly",
			line:    []byte("apple"),
			pattern: "^apple$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "^apple$ does not match apple pie",
			line:    []byte("apple pie"),
			pattern: "^apple$",
			want:    false,
			wantErr: false,
		},
		{
			name:    "^apple$ does not match green apple",
			line:    []byte("green apple"),
			pattern: "^apple$",
			want:    false,
			wantErr: false,
		},
		{
			name:    "^cat$ matches cat exactly",
			line:    []byte("cat"),
			pattern: "^cat$",
			want:    true,
			wantErr: false,
		},
		{
			name:    "^cat$ does not match cats",
			line:    []byte("cats"),
			pattern: "^cat$",
			want:    false,
			wantErr: false,
		},
		// + quantifier tests (one or more)
		{
			name:    "ca+ts matches cats (1 a)",
			line:    []byte("cats"),
			pattern: "ca+ts",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ca+ts matches caats (2 a's)",
			line:    []byte("caats"),
			pattern: "ca+ts",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ca+ts matches caaats (3 a's)",
			line:    []byte("caaats"),
			pattern: "ca+ts",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ca+ts does not match cts (0 a's)",
			line:    []byte("cts"),
			pattern: "ca+ts",
			want:    false,
			wantErr: false,
		},
		{
			name:    "a+ matches a",
			line:    []byte("a"),
			pattern: "a+",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+ matches aaa",
			line:    []byte("aaa"),
			pattern: "a+",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+ does not match b",
			line:    []byte("b"),
			pattern: "a+",
			want:    false,
			wantErr: false,
		},
		{
			name:    "a+b matches ab",
			line:    []byte("ab"),
			pattern: "a+b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+b matches aaab",
			line:    []byte("aaab"),
			pattern: "a+b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+b does not match b",
			line:    []byte("b"),
			pattern: "a+b",
			want:    false,
			wantErr: false,
		},
		{
			name:    "a+b+ matches ab",
			line:    []byte("ab"),
			pattern: "a+b+",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+b+ matches aaabbb",
			line:    []byte("aaabbb"),
			pattern: "a+b+",
			want:    true,
			wantErr: false,
		},
		// Backtracking tests (cases where greedy matching would fail)
		{
			name:    "ca+at matches caaats (backtracking needed)",
			line:    []byte("caaats"),
			pattern: "ca+at",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ca+at matches caats (no backtracking needed)",
			line:    []byte("caats"),
			pattern: "ca+at",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ca+at does not match cts",
			line:    []byte("cts"),
			pattern: "ca+at",
			want:    false,
			wantErr: false,
		},
		{
			name:    "ca+at does not match caaaa (no 't' at end)",
			line:    []byte("caaaa"),
			pattern: "ca+at",
			want:    false,
			wantErr: false,
		},
		{
			name:    "a+b matches aaab (backtracking needed)",
			line:    []byte("aaab"),
			pattern: "a+b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d+apple matches 123apple (backtracking might be needed)",
			line:    []byte("123apple"),
			pattern: "\\d+apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\d+ apple matches 100 apples (backtracking for space)",
			line:    []byte("100 apples"),
			pattern: "\\d+ apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+b+c matches aaabbbc (multiple quantifiers with backtracking)",
			line:    []byte("aaabbbc"),
			pattern: "a+b+c",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a+b+c matches abc (minimal match)",
			line:    []byte("abc"),
			pattern: "a+b+c",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc]+d matches abcd (backtracking with char class)",
			line:    []byte("abcd"),
			pattern: "[abc]+d",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc]+d matches aaabbbcccd (backtracking with char class)",
			line:    []byte("aaabbbcccd"),
			pattern: "[abc]+d",
			want:    true,
			wantErr: false,
		},
		// Literal backslash tests
		{
			name:    "\\\\ matches backslash character",
			line:    []byte("a\\b"),
			pattern: "a\\\\b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\\\ does not match when no backslash present",
			line:    []byte("ab"),
			pattern: "a\\\\b",
			want:    false,
			wantErr: false,
		},
		{
			name:    "\\\\d matches backslash followed by 'd'",
			line:    []byte("test\\d123"),
			pattern: "\\\\d",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\\\d does not match digit (should match literal \\d)",
			line:    []byte("test5"),
			pattern: "\\\\d",
			want:    false,
			wantErr: false,
		},
		// Zero or One (?) quantifier tests
		{
			name:    "dogs? matches dog (0 s)",
			line:    []byte("dog"),
			pattern: "dogs?",
			want:    true,
			wantErr: false,
		},
		{
			name:    "dogs? matches dogs (1 s)",
			line:    []byte("dogs"),
			pattern: "dogs?",
			want:    true,
			wantErr: false,
		},
		{
			name:    "colou?r matches color (0 u)",
			line:    []byte("color"),
			pattern: "colou?r",
			want:    true,
			wantErr: false,
		},
		{
			name:    "colou?r matches colour (1 u)",
			line:    []byte("colour"),
			pattern: "colou?r",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a? matches empty at start",
			line:    []byte("b"),
			pattern: "a?b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a? matches single a",
			line:    []byte("ab"),
			pattern: "a?b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\\\d? matches 0 digits",
			line:    []byte("apple"),
			pattern: "\\d?apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "\\\\d? matches 1 digit",
			line:    []byte("1apple"),
			pattern: "\\d?apple",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc]? matches 0 chars",
			line:    []byte("dog"),
			pattern: "[abc]?dog",
			want:    true,
			wantErr: false,
		},
		{
			name:    "[abc]? matches 1 char",
			line:    []byte("adog"),
			pattern: "[abc]?dog",
			want:    true,
			wantErr: false,
		},
		// Dot (.) metacharacter tests
		{
			name:    "d.g matches dog",
			line:    []byte("dog"),
			pattern: "d.g",
			want:    true,
			wantErr: false,
		},
		{
			name:    "d.g matches dag",
			line:    []byte("dag"),
			pattern: "d.g",
			want:    true,
			wantErr: false,
		},
		{
			name:    "d.g matches d1g",
			line:    []byte("d1g"),
			pattern: "d.g",
			want:    true,
			wantErr: false,
		},
		{
			name:    "d.g matches d@g",
			line:    []byte("d@g"),
			pattern: "d.g",
			want:    true,
			wantErr: false,
		},
		{
			name:    "d.g does not match cog (first char mismatch)",
			line:    []byte("cog"),
			pattern: "d.g",
			want:    false,
			wantErr: false,
		},
		{
			name:    "d.g does not match dg (missing middle char)",
			line:    []byte("dg"),
			pattern: "d.g",
			want:    false,
			wantErr: false,
		},
		{
			name:    "... matches any 3 chars",
			line:    []byte("abc"),
			pattern: "...",
			want:    true,
			wantErr: false,
		},
		{
			name:    "... matches 123",
			line:    []byte("123"),
			pattern: "...",
			want:    true,
			wantErr: false,
		},
		{
			name:    ".+ matches one or more chars",
			line:    []byte("hello"),
			pattern: ".+",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a.?b matches ab (0 chars)",
			line:    []byte("ab"),
			pattern: "a.?b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "a.?b matches acb (1 char)",
			line:    []byte("acb"),
			pattern: "a.?b",
			want:    true,
			wantErr: false,
		},
		{
			name:    "c.t matches cat",
			line:    []byte("cat"),
			pattern: "c.t",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchLine(tt.line, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
