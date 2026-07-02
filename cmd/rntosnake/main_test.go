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
	}{
		{"strings2 basic", "helloWorld", "hello_world"},
		{"strings2 space handling", "this is a test", "this_is_a_test"},
		{"strings2 dash conversion", "hello-world", "hello_world"},
		{"strings2 double caps", "HTTPResponse", "http_response"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if strings.ToLower(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
