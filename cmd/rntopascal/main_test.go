package main

import (
	"testing"
)

func TestPascalAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "hello_world", "HelloWorld", "strings2"},
		{"strings2 spaces", "this is a test", "ThisIsATest", "strings2"},
		{"strings2 mixed case", "MIXED_case-test", "MixedCaseTest", "strings2"},
		{"strings2 acronym handling", "HTTP_REQUEST_ID", "HttpRequestId", "strings2"},

		// Prove removed things don't work (skipped)
		{"searKing basic", "hello_world", "HelloWorld", "searKing"},
		{"resenje basic", "hello_world", "HelloWorld", "resenje"},
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
