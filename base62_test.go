package main

import "testing"

func TestEncodeBase62(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{0, "a"},
		{1, "b"},
		{61, "9"},
		{62, "ba"},
		{12345, "dnh"},
	}

	for _, test := range tests {
		result := encodeBase62(test.input)
		if result != test.expected {
			t.Errorf("encodeBase62(%d) = %s, want %s", test.input, result, test.expected)
		}
	}
}

func TestEncodeBase62_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := int64(1); i <= 10000; i++ {
		code := encodeBase62(i)
		if seen[code] {
			t.Fatalf("duplicate code %s for id %d", code, i)
		}
		seen[code] = true
	}
}
