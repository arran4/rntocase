package main

import (
	"github.com/iancoleman/strcase"
	"testing"
)

func TestSnakeAlgorithms(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", "helloWorld"},
		{"space handling", "this is a test"},
		{"dash conversion", "hello-world"},
		{"double caps", "HTTPResponse"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strcase.ToSnake(tt.input)
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
