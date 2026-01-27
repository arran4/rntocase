package rntocase

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func BenchmarkRenameFiles_Unbuffered(b *testing.B) {
	// Generate dummy files
	files := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = fmt.Sprintf("path/to/file_%d.txt", i)
	}

	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	// Redirect stdout to dev/null to simulate I/O without printing to console
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("failed to open devnull: %v", err)
	}
	defer f.Close()

	oldStdout := os.Stdout
	os.Stdout = f
	defer func() { os.Stdout = oldStdout }()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Run in dry-run mode
		_ = RenameFiles(os.Stdout, files, renameFunc, true, false)
	}
}

func BenchmarkRenameFiles_Buffered(b *testing.B) {
	// Generate dummy files
	files := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = fmt.Sprintf("path/to/file_%d.txt", i)
	}

	renameFunc := func(s string) (string, error) {
		return s + "_renamed", nil
	}

	// Redirect stdout to dev/null
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0666)
	if err != nil {
		b.Fatalf("failed to open devnull: %v", err)
	}
	defer f.Close()

	buf := bufio.NewWriter(f)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Run in dry-run mode
		_ = RenameFiles(buf, files, renameFunc, true, false)
		buf.Flush()
	}
}
