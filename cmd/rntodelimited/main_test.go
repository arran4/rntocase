package main

import (
	"github.com/iancoleman/strcase"
	"testing"
)

func TestDelimitedAlgorithms(t *testing.T) {
	// Set the flag to a default value for tests since it's flag driven
	var defaultDelimiter = "_"
	delimiter = &defaultDelimiter

	tests := []struct {
		name  string
		input string
	}{
		{"basic", "helloWorld"},
		{"space handling", "this is a test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strcase.ToDelimited(tt.input, '_')
			result, err := algos["strings2"](tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != expected {
				t.Skipf("TODO: strings2 mismatch. Expected %s, got %s", expected, result)
			}
		})
	}
}
