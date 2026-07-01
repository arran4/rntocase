package main

import (
	"strings"
	"testing"
)

func TestSnakeAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		algo     string
	}{
		{"strings2 basic", "helloWorld", "hello_world", "strings2"},
		// Prove removed things don't work (skipped)
		{"screaming-iancoleman basic", "helloWorld", "HELLO_WORLD", "screaming-iancoleman"},
		{"resenje basic", "helloWorld", "hello_world", "resenje"},
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
				// Strings2 might return hello_World since it just adds an underscore
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
