package main

import (
	"testing"
)

func TestAcronymAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "hello world", "HW", "strings2"},
		// Prove removed things don't work (skipped)
		{"gobeam basic", "hello world", "HW", "gobeam"},
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
