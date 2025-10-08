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
