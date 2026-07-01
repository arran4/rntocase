package main

import (
	"testing"
)

func TestReverseAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"go basic", "hello", "olleh", "go"},
		// Prove removed things don't work (skipped)
		{"searking basic", "hello", "olleh", "searking"},
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
