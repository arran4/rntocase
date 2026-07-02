package main

import (
	"testing"
)

func TestCamelAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strings2 basic", "hello_world", "helloWorld"},
		{"strings2 spaces", "this is a test", "thisIsATest"},
		{"strings2 mixed case", "MIXED_case-test", "mixedCaseTest"},
		{"strings2 acronym handling", "HTTP_REQUEST_ID", "httpRequestId"},
		{"strings2 leading underscore", "_hello_world_", "helloWorld"},
		{"strings2 double delimiters", "hello__world--test", "helloWorldTest"},
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
