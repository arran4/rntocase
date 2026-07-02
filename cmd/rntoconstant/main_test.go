package main

import (
	"testing"
)

func TestConstantAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strings2 basic", "helloWorld", "HELLO_WORLD"},
		{"strings2 spaces", "this is a test", "THIS_IS_A_TEST"},
		{"strings2 acronym handling", "httpRequest", "HTTP_REQUEST"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
