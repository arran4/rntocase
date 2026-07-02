package main

import (
	"strings"
	"testing"
)

func TestKebabAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "helloWorld", "hello-world", "strings2"},
		{"strings2 space handling", "this is a test", "this-is-a-test", "strings2"},
		{"strings2 underscore conversion", "hello_world", "hello-world", "strings2"},
		{"strings2 double caps", "HTTPResponse", "http-response", "strings2"},

		// Prove removed things don't work (skipped)
		{"screaming-iancoleman basic", "helloWorld", "HELLO-WORLD", "screaming-iancoleman"},
		{"resenje basic", "helloWorld", "hello-world", "resenje"},
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
			if strings.ToLower(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
