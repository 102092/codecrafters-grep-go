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
		// 에러 케이스
		{
			name:    "multi-char pattern error",
			line:    []byte("test"),
			pattern: "abc",
			want:    false,
			wantErr: true,
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
