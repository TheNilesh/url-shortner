package svc

import (
	"testing"
)

func TestRandomStrGen_Generate(t *testing.T) {
	tests := []struct {
		name      string
		minLength int
		maxLength int
		charset   string
	}{
		{
			name:      "generates string of minimum length",
			minLength: 5,
			maxLength: 10,
			charset:   "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "generates string of maximum length",
			minLength: 5,
			maxLength: 5,
			charset:   "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name:      "generates string of random length",
			minLength: 5,
			maxLength: 10,
			charset:   "abcdefghijklmnopqrstuvwxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRandomStrGen(tt.minLength, tt.maxLength, tt.charset)
			got := r.Generate()
			if len(got) < tt.minLength || len(got) > tt.maxLength {
				t.Errorf("Generate() = %q, want string of length between %d and %d", got, tt.minLength, tt.maxLength)
			}
			for _, c := range got {
				if !contains(tt.charset, string(c)) {
					t.Errorf("Generate() = %q, contains invalid character %q", got, string(c))
				}
			}
		})
	}
}

func contains(s string, c string) bool {
	for _, r := range s {
		if string(r) == c {
			return true
		}
	}
	return false
}
