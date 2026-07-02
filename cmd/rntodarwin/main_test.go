package main

import (
	"testing"
)

func TestDarwinAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strings2 basic", "hello world", "Hello_World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
