package main

import (
	"github.com/iancoleman/strcase"
	"testing"
)

func TestCamelAlgorithms(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", "hello_world"},
		{"spaces", "this is a test"},
		{"mixed case", "MIXED_case-test"},
		{"acronym handling", "HTTP_REQUEST_ID"},
		{"leading underscore", "_hello_world_"},
		{"double delimiters", "hello__world--test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strcase.ToCamel(tt.input)
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
