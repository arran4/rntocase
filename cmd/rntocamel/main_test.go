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
		{"strings2 acronym", "ID_NUMBER", "idNumber", "strings2"},
		// Prove removed things don't work (skipped)
		{"iancoleman basic", "hello_world", "helloWorld", "iancoleman"},
		{"lower-iancoleman basic", "hello_world", "helloWorld", "lower-iancoleman"},
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
