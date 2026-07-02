package main

import (
	"testing"
)

func TestPascalAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strings2 basic", "hello_world", "HelloWorld"},
		{"strings2 spaces", "this is a test", "ThisIsATest"},
		{"strings2 mixed case", "MIXED_case-test", "MixedCaseTest"},
		{"strings2 acronym handling", "HTTP_REQUEST_ID", "HttpRequestId"},
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
