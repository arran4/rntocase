package main

import (
	"strings"
	"testing"
)

func TestDelimitedAlgorithms(t *testing.T) {
	// Set the flag to a default value for tests since it's flag driven
    var defaultDelimiter = "_"
	delimiter = &defaultDelimiter

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strings2 basic", "helloWorld", "hello_world"},
		{"strings2 space handling", "this is a test", "this_is_a_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := algos["strings2"](tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if strings.ToLower(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
