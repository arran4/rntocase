package main

import (
	"github.com/iancoleman/strcase"
	"testing"
)

func TestConstantAlgorithms(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", "helloWorld"},
		{"spaces", "this is a test"},
		{"acronym handling", "httpRequest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strcase.ToScreamingSnake(tt.input)
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != expected {
				t.Errorf("TODO: strings2 mismatch. Expected %s, got %s", expected, result)
			}
		})
	}
}
