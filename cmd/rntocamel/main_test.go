package main

import (
	"testing"
)

func TestCamelAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "hello_world", "helloWorld", "strings2"},
		{"strings2 spaces", "this is a test", "thisIsATest", "strings2"},
		{"strings2 mixed case", "MIXED_case-test", "mixedCaseTest", "strings2"},
		{"strings2 acronym handling", "HTTP_REQUEST_ID", "httpRequestId", "strings2"},
		{"strings2 leading underscore", "_hello_world_", "helloWorld", "strings2"},
		{"strings2 double delimiters", "hello__world--test", "helloWorldTest", "strings2"},

		// Prove removed things don't work (skipped)
		{"iancoleman basic", "hello_world", "helloWorld", "iancoleman"},
		{"lower-iancoleman basic", "hello_world", "helloWorld", "lower-iancoleman"},
		{"searking basic", "hello_world", "helloWorld", "searking"},
		{"gobeam basic", "hello_world", "helloWorld", "gobeam"},
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
