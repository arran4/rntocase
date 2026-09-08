package cli

import (
	"github.com/arran4/strings2"
	"github.com/iancoleman/strcase"
	"strings"
	"testing"
)

func TestCamelAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "hello_world", "helloWorld"},
		{"spaces", "this is a test", "thisIsATest"},
		{"mixed case", "MIXED_case-test", "mixedCaseTest"},
		{"acronym handling", "HTTP_REQUEST_ID", "httpRequestId"},
		{"leading underscore", "_hello_world_", "helloWorld"},
		{"double delimiters", "hello__world--test", "helloWorldTest"},
	}

	converter := func(s string) (string, error) {
		return strings2.ToCamel(s, strings2.ParserEmitEmpty(true))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("strings2 mismatch. Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestPascalAlgorithms(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"basic", "hello_world", "HelloWorld"},
		{"spaces", "this is a test", "ThisIsATest"},
		{"mixed case", "MIXED_case-test", "MixedCaseTest"},
		{"acronym handling", "HTTP_REQUEST_ID", "HttpRequestId"},
		{"leading underscore", "_hello_world_", "HelloWorld"},
		{"double delimiters", "hello__world--test", "HelloWorldTest"},
	}

	converter := func(s string) (string, error) {
		return strings2.ToPascal(s, strings2.ParserEmitEmpty(true))
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("strings2 mismatch. Expected %s, got %s", tt.expected, result)
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

func TestTrimAlgorithms(t *testing.T) {
	tests := []struct {
		name      string
		trimChars string
		input     string
		expected  string
	}{
		{"spaces default", "", "   hello world   ", "hello world"},
		{"trim underscore", "_", "_hello_world_", "hello_world"},
		{"trim dash", "-", "-hello-world-", "hello-world"},
		{"trim multiple chars", "_-", "_-hello-_world-_", "hello-_world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := func(s string) (string, error) {
				if tt.trimChars == "" {
					return strings.TrimSpace(s), nil
				} else {
					return strings.Trim(s, tt.trimChars), nil
				}
			}
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Trim mismatch. Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestDotAlgorithms(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		input     string
		expected  string
	}{
		{"default delimiter", "", "hello world", "hello.world"},
		{"custom delimiter", "-", "hello world", "hello-world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delim := tt.delimiter
			if delim == "" {
				delim = "."
			}
			converter := func(s string) (string, error) {
				return strings2.ToFormattedString(s, strings2.OptionDelimiter(delim), strings2.OptionFirstLower(), strings2.ParserEmitEmpty(true))
			}
			result, err := converter(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Dot mismatch. Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
