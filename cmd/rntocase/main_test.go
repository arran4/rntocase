package main

import (
	"github.com/iancoleman/strcase"
	"github.com/arran4/strings2"
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

	converter := func(s string) (string, error) {
		return strings2.ToPascal(s, strings2.ParserEmitEmpty(true))
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

func TestSnakeAlgorithms(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", "HelloWorld"},
		{"spaces", "this is a test"},
		{"mixed case", "MIXED_case-test"},
		{"acronym handling", "HTTPRequestID"},
		{"leading underscore", "_hello_world_"},
		{"double delimiters", "hello__world--test"},
	}

	converter := func(s string) (string, error) {
		return strings2.ToSnake(s, strings2.OptionLoose(), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
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

func TestDelimitedAlgorithms(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"basic", "HelloWorld"},
		{"spaces", "this is a test"},
		{"mixed case", "MIXED_case-test"},
		{"acronym handling", "HTTPRequestID"},
		{"leading underscore", "_hello_world_"},
		{"double delimiters", "hello__world--test"},
	}

	delimiter := "."
	ignore := ""
	converter := func(s string) (string, error) {
		return strings2.ToFormattedString(s, strings2.OptionDelimiter(delimiter), strings2.OptionIgnore(ignore), strings2.OptionCaseMode(strings2.CMWhispering), strings2.ParserEmitEmpty(true))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := strcase.ToDelimited(tt.input, '.')
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
