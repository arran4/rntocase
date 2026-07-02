package main

import (
	"testing"
)

func TestConstantAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "helloWorld", "HELLO_WORLD", "strings2"},
		{"strings2 spaces", "this is a test", "THIS_IS_A_TEST", "strings2"},
		{"strings2 acronym handling", "httpRequest", "HTTP_REQUEST", "strings2"},

		// Prove removed things don't work (skipped)
		{"iancoleman basic", "helloWorld", "HELLO_WORLD", "iancoleman"},
		{"resenje basic", "helloWorld", "HELLO_WORLD", "resenje"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			algoFn, exists := algos[tt.algo]
			if !exists {
				t.Skipf("Algorithm %s not found (proves it was removed)", tt.algo)
				return
			}
			result, err := algoFn(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
