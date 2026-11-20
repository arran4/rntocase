package cli

import (
	"strings"
	"testing"
)

func TestRenameActionsRequireFileOperands(t *testing.T) {
	tests := []struct {
		name string
		run  func() error
	}{
		{name: "acronym", run: func() error { return RunAcronym(false, false, false, false, nil, nil) }},
		{name: "camel", run: func() error { return RunCamel(false, false, false, false, nil, nil) }},
		{name: "constant", run: func() error { return RunConstant(false, false, false, false, nil, nil) }},
		{name: "darwin", run: func() error { return RunDarwin(false, false, false, false, nil, nil) }},
		{name: "delimited", run: func() error { return RunDelimited("", "", false, false, false, false, nil, nil) }},
		{name: "dot", run: func() error { return RunDot("", false, false, false, false, nil, nil) }},
		{name: "kebab", run: func() error { return RunKebab(false, false, false, false, nil, nil) }},
		{name: "lower", run: func() error { return RunLower(false, false, false, false, nil, nil) }},
		{name: "lower-leading", run: func() error { return RunLowerLeading(false, false, false, false, nil, nil) }},
		{name: "pascal", run: func() error { return RunPascal(false, false, false, false, nil, nil) }},
		{name: "reverse", run: func() error { return RunReverse(false, false, false, false, false, nil, nil) }},
		{name: "snake", run: func() error { return RunSnake(false, false, false, false, nil, nil) }},
		{name: "title", run: func() error { return RunTitle(false, false, false, false, nil, nil) }},
		{name: "upper", run: func() error { return RunUpper(false, false, false, false, nil, nil) }},
		{name: "upper-leading", run: func() error { return RunUpperLeading(false, false, false, false, nil, nil) }},
		{name: "trim", run: func() error { return RunTrim("", false, false, false, false, nil, nil) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if err == nil {
				t.Fatal("expected error for missing file operands")
			}
			if !strings.Contains(err.Error(), "rntocase "+tt.name+":") ||
				!strings.Contains(err.Error(), "at least one file or directory") {
				t.Fatalf("unexpected missing-input error: %v", err)
			}
		})
	}
}

func TestRequireRenameFileOperandsAllowsProvidedFiles(t *testing.T) {
	if err := requireRenameFileOperands("test", []string{"example.txt"}); err != nil {
		t.Fatalf("unexpected error for provided filename: %v", err)
	}
}
