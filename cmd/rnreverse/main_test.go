package main

import (
	"testing"
)

func TestReverseAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
        mode     string
	}{
		{"go basic", "hello", "olleh", "reverse_runes"},
        {"word reverse", "hello world test", "test world hello", "reverse_words"},
	}

    algs := map[string]func(string) (string, error){
        "reverse_runes": converter,
        "reverse_words": wordReverseConverter,
    }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algs[tt.mode](tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
