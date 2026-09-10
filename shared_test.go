package rntocase

import (
	"encoding/json"
	"fmt"
	"io"

	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSplitExtension(t *testing.T) {
	cases := []struct {
		filename string
		wantName string
		wantExt  string
	}{
		{"hello.txt", "hello", ".txt"},
		{"noextension", "noextension", ""},
		{".env", ".env", ""},
		{".gitignore", ".gitignore", ""},
		{"archive.tar.gz", "archive", ".tar.gz"},
		{"archive.backup.tar.gz", "archive.backup", ".tar.gz"},
		{".archive.tar.gz", ".archive", ".tar.gz"},
		{".tar.gz", ".tar.gz", ""},
		{".dot.file", ".dot", ".file"},
		{"", "", ""},
		{".", ".", ""},
		{"..", "..", ""},
	}

	for _, tc := range cases {
		name, ext := SplitExtension(tc.filename)
		if name != tc.wantName || ext != tc.wantExt {
			t.Errorf("SplitExtension(%q) = %q, %q; want %q, %q", tc.filename, name, ext, tc.wantName, tc.wantExt)
		}
	}
}

func TestRenameFiles(t *testing.T) {
	// Rename function (to lowercase)
	renameFunc := func(s string) (string, error) {
		return strings.ToLower(s), nil
	}

	t.Run("successful rename and extension preservation", func(t *testing.T) {
		tempDir := t.TempDir()

		files := []string{
			"Hello World.txt",
			"testFile.go",
			"Archive.tar.gz",
			".Env",
		}

		var paths []string
		for _, f := range files {
			path := filepath.Join(tempDir, f)
			if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
				t.Fatalf("Failed to create file %s: %v", path, err)
			}
			paths = append(paths, path)
		}

		if err := RenameFiles(paths, renameFunc, false, false, false); err != nil {
			t.Fatalf("RenameFiles failed: %v", err)
		}

		expectedFiles := []string{
			"hello world.txt",
			"testfile.go",
			"archive.tar.gz",
			".env",
		}

		for _, f := range expectedFiles {
			path := filepath.Join(tempDir, f)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected file %s does not exist", path)
			}
		}

		// Ensure originals are gone where they changed case (if FS is case-sensitive, otherwise it just renamed)
		// For robustness across OS, we check the new files exist and we successfully returned nil.
	})

	t.Run("collision existing unrelated destination", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "A.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile(filepath.Join(tempDir, "a.txt"), []byte("content"), 0644)
		// If we are on a case-insensitive FS, A.txt and a.txt might be the same.
		// Let's use totally different names to be safe.

		// Safer approach for cross-platform collision test
		if err := os.WriteFile(filepath.Join(tempDir, "SourceFile.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile(filepath.Join(tempDir, "sourcefile.txt"), []byte("different"), 0644)
		// For case-insensitive FS, it will just overwrite content, let's use different base names.

		// Unrelated collision test:
		if err := os.WriteFile(filepath.Join(tempDir, "Foo.txt"), []byte("foo"), 0644); err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile(filepath.Join(tempDir, "foo.txt"), []byte("existing foo"), 0644)
		// On Case-Insensitive FS (macOS/Windows) this might just modify Foo.txt.
		// We can test collision by mapping 'Bar.txt' to 'foo.txt' by just a dumb replace func.

		if err := os.WriteFile(filepath.Join(tempDir, "Bar.txt"), []byte("bar"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "baz.txt"), []byte("baz"), 0644); err != nil {
			t.Fatal(err)
		}

		renameToBaz := func(s string) (string, error) {
			return "baz", nil
		}

		err := RenameFiles([]string{filepath.Join(tempDir, "Bar.txt")}, renameToBaz, false, false, false)
		if err == nil {
			t.Fatal("Expected collision error, got nil")
		}
		if !strings.Contains(err.Error(), "collision: destination") {
			t.Errorf("Expected collision error message, got: %v", err)
		}
	})

	t.Run("two sources map to same destination", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "File1.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "File2.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		renameToSame := func(s string) (string, error) {
			return "same", nil
		}

		paths := []string{filepath.Join(tempDir, "File1.txt"), filepath.Join(tempDir, "File2.txt")}
		err := RenameFiles(paths, renameToSame, false, false, false)
		if err == nil {
			t.Fatal("Expected error due to multiple files mapping to same destination")
		}
		if !strings.Contains(err.Error(), "multiple source files map to destination") {
			t.Errorf("Expected specific collision error, got: %v", err)
		}
	})

	t.Run("two sources map to same destination with mixed relative absolute paths", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "FOO.txt"), []byte("FIRST"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "BAR.txt"), []byte("SECOND"), 0644); err != nil {
			t.Fatal(err)
		}

		renameToLower := func(s string) (string, error) {
			if s == "BAR" {
				return "foo", nil
			}
			return strings.ToLower(s), nil
		}

		// Change directory to tempDir to test true relative paths against absolute ones
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(tempDir); err != nil {
			t.Fatal(err)
		}
		defer func() {
			_ = os.Chdir(cwd)
		}()

		// Use a true relative path ("FOO.txt") and an absolute path
		absPath, _ := filepath.Abs(".")
		paths := []string{
			"FOO.txt",
			filepath.Join(absPath, "BAR.txt"),
		}

		err = RenameFiles(paths, renameToLower, false, false, false)
		if err == nil {
			t.Fatal("Expected error due to multiple files mapping to same destination across rel/abs paths")
		}

		if !strings.Contains(err.Error(), "multiple source files map to destination") {
			t.Errorf("Expected collision error message, got: %v", err)
		}

		// Ensure filesystem unchanged
		b1, _ := os.ReadFile("FOO.txt")
		if string(b1) != "FIRST" {
			t.Errorf("Source contents changed, expected FIRST, got %s", string(b1))
		}
		b2, _ := os.ReadFile(filepath.Join(absPath, "BAR.txt"))
		if string(b2) != "SECOND" {
			t.Errorf("Source contents changed, expected SECOND, got %s", string(b2))
		}
	})

	t.Run("case insensitive destination collision", func(t *testing.T) {
		tempDir := t.TempDir()

		// Write two files that will be trimmed to destinations that differ only in case.
		if err := os.WriteFile(filepath.Join(tempDir, " foo.TXT"), []byte("FIRST"), 0644); err != nil {
			t.Fatal(err)
		}
		// Write to "bar .txt" which we'll trim to "bar.txt" and rename to "foo.txt" to avoid
		// initial creation overwriting on case-insensitive filesystems.
		if err := os.WriteFile(filepath.Join(tempDir, " bar.txt"), []byte("SECOND"), 0644); err != nil {
			t.Fatal(err)
		}

		// A trim function that mimics trimming space and returning the name
		trimFunc := func(s string) (string, error) {
			trimmed := strings.TrimSpace(s)
			if trimmed == "bar" {
				return "foo", nil // simulate mapping to same base name but with .txt instead of .TXT
			}
			return trimmed, nil
		}

		paths := []string{
			filepath.Join(tempDir, " foo.TXT"),
			filepath.Join(tempDir, " bar.txt"),
		}

		err := RenameFiles(paths, trimFunc, false, false, false)
		if err == nil {
			t.Fatal("Expected error due to case-insensitive multiple files mapping to same destination")
		}

		if !strings.Contains(err.Error(), "multiple source files map to destination") {
			t.Errorf("Expected collision error message, got: %v", err)
		}

		// Ensure filesystem unchanged
		b1, _ := os.ReadFile(filepath.Join(tempDir, " foo.TXT"))
		if string(b1) != "FIRST" {
			t.Errorf("Source contents changed, expected FIRST, got %s", string(b1))
		}
		b2, _ := os.ReadFile(filepath.Join(tempDir, " bar.txt"))
		if string(b2) != "SECOND" {
			t.Errorf("Source contents changed, expected SECOND, got %s", string(b2))
		}
	})

	t.Run("inspect directory entries without following symlinks", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "Foo.txt"), []byte("FIRST"), 0644); err != nil {
			t.Fatal(err)
		}

		// Create a dangling symlink `foo.txt -> missing.txt`
		danglingPath := filepath.Join(tempDir, "foo.txt")
		if err := os.Symlink("missing.txt", danglingPath); err != nil {
			t.Skipf("Skipping symlink test: %v", err)
		}

		renameToLower := func(s string) (string, error) {
			return strings.ToLower(s), nil
		}

		paths := []string{filepath.Join(tempDir, "Foo.txt")}

		// Should error because `foo.txt` exists as a symlink
		err := RenameFiles(paths, renameToLower, false, false, false)
		if err == nil {
			t.Fatal("Expected collision error with dangling symlink, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Expected already exists error message, got: %v", err)
		}

		// Verify symlink still there and intact
		destStat, destErr := os.Lstat(danglingPath)
		if destErr != nil {
			t.Errorf("Expected dangling symlink to remain, got err: %v", destErr)
		} else if destStat.Mode()&os.ModeSymlink == 0 {
			t.Errorf("Expected dangling path to still be a symlink")
		}
	})

	t.Run("inspect directory entries symlink pointing to source", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "Foo.txt"), []byte("FIRST"), 0644); err != nil {
			t.Fatal(err)
		}

		// Create a symlink `foo.txt -> Foo.txt`
		symlinkPath := filepath.Join(tempDir, "foo.txt")
		if err := os.Symlink("Foo.txt", symlinkPath); err != nil {
			t.Skipf("Skipping symlink test: %v", err)
		}

		renameToLower := func(s string) (string, error) {
			return strings.ToLower(s), nil
		}

		paths := []string{filepath.Join(tempDir, "Foo.txt")}

		// Should error because `foo.txt` exists as a symlink
		err := RenameFiles(paths, renameToLower, false, false, false)
		if err == nil {
			t.Fatal("Expected collision error with symlink pointing to source, got nil")
		}

		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("Expected already exists error message, got: %v", err)
		}

		// Verify symlink still there and intact
		destStat, destErr := os.Lstat(symlinkPath)
		if destErr != nil {
			t.Errorf("Expected symlink to remain, got err: %v", destErr)
		} else if destStat.Mode()&os.ModeSymlink == 0 {
			t.Errorf("Expected symlink path to still be a symlink")
		}

		// Verify symlink target is correct
		target, err := os.Readlink(symlinkPath)
		if err != nil {
			t.Errorf("Failed to readlink: %v", err)
		}
		if target != "Foo.txt" {
			t.Errorf("Expected symlink to point to Foo.txt, got %s", target)
		}
	})

	t.Run("dry run collision detection", func(t *testing.T) {
		tempDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(tempDir, "File1.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "File2.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		renameToSame := func(s string) (string, error) {
			return "same", nil
		}

		paths := []string{filepath.Join(tempDir, "File1.txt"), filepath.Join(tempDir, "File2.txt")}
		err := RenameFiles(paths, renameToSame, true, false, false)
		if err == nil {
			t.Fatal("Expected dry-run to still detect planning collisions")
		}
	})

	t.Run("mixed success and failure during execution", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "Good.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		// To simulate an execution failure, we'll try to rename a file that we delete before os.Rename,
		// or on unix, remove write permissions from directory. Here, we can just delete it after planning.
		if err := os.WriteFile(filepath.Join(tempDir, "Bad1.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		renameFuncLower := func(s string) (string, error) {
			return strings.ToLower(s), nil
		}

		// We will do a bit of a hack: to simulate a syscall failure during execution,
		// we pass the files to RenameFiles but we run a goroutine that deletes Bad1.txt right after planning?
		// Actually, simpler: just pass a file that doesn't exist, but bypass stat check for existing file?
		// Wait, if it doesn't exist, the stat for source will fail but we only check destination.
		// Let's test that: we provide a non-existent file as source.
		paths := []string{
			filepath.Join(tempDir, "Good.txt"),
			filepath.Join(tempDir, "Bad1.txt"),
		}

		// Remove it so the syscall fails
		_ = os.Remove(filepath.Join(tempDir, "Bad1.txt"))

		err := RenameFiles(paths, renameFuncLower, false, false, false)
		if err == nil {
			t.Fatal("Expected error for the bad files")
		}

		// Good.txt should be renamed because execution continues
		if _, err := os.Stat(filepath.Join(tempDir, "good.txt")); os.IsNotExist(err) {
			t.Errorf("Expected good.txt to exist, it was not renamed")
		}
	})

	t.Run("batch aborted on planning collision", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "Good.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "Bad1.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(tempDir, "Bad2.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		renameMixed := func(s string) (string, error) {
			if strings.HasPrefix(s, "Bad") {
				return "same_bad", nil // Will cause collision
			}
			return strings.ToLower(s), nil
		}

		paths := []string{
			filepath.Join(tempDir, "Good.txt"),
			filepath.Join(tempDir, "Bad1.txt"),
			filepath.Join(tempDir, "Bad2.txt"),
		}
		err := RenameFiles(paths, renameMixed, false, false, false)
		if err == nil {
			t.Fatal("Expected error for the bad files")
		}

		// Good.txt should NOT be renamed because the batch should be aborted
		if _, err := os.Stat(filepath.Join(tempDir, "Good.txt")); os.IsNotExist(err) {
			t.Errorf("Expected Good.txt to still exist, but batch was not aborted")
		}
	})

	t.Run("unchanged path skips rename", func(t *testing.T) {
		tempDir := t.TempDir()

		if err := os.WriteFile(filepath.Join(tempDir, "alreadylower.txt"), []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}

		paths := []string{filepath.Join(tempDir, "alreadylower.txt")}
		err := RenameFiles(paths, renameFunc, false, false, false)
		if err != nil {
			t.Fatalf("Expected nil error for unchanged, got: %v", err)
		}
	})
}

func TestConfirm(t *testing.T) {
	// Create a pipe to simulate stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	reader := bufio.NewReader(r)

	// Write two responses to the pipe: "y\n" and "y\n"
	go func() {
		defer func() {
			_ = w.Close()
		}()
		_, err := io.WriteString(w, "y\ny\n")
		if err != nil {
			t.Errorf("Failed to write to pipe: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
	}()

	if !ConfirmWithReader("Confirm 1?", reader) {
		t.Error("Failed to confirm 1")
	}

	if !ConfirmWithReader("Confirm 2?", reader) {
		t.Error("Failed to confirm 2")
	}
}

func TestRenameFilesJSON(t *testing.T) {
	tempDir := t.TempDir()

	// Helper to capture stdout
	captureStdout := func(f func()) string {
		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("Failed to create pipe: %v", err)
		}
		os.Stdout = w

		errCh := make(chan error, 2)

		go func() {
			f()
			errCh <- w.Close()
		}()

		var buf strings.Builder
		_, copyErr := io.Copy(&buf, r)
		errCh <- copyErr

		os.Stdout = oldStdout

		closeErr1 := <-errCh
		closeErr2 := <-errCh
		if closeErr1 != nil {
			t.Fatalf("Error in capture goroutine: %v", closeErr1)
		}
		if closeErr2 != nil {
			t.Fatalf("Error copying from pipe: %v", closeErr2)
		}

		if err := r.Close(); err != nil {
			t.Fatalf("Error closing pipe reader: %v", err)
		}

		return buf.String()
	}

	renameFunc := func(s string) (string, error) {
		return strings.ToUpper(s), nil
	}

	t.Run("successful dry run JSON", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "dry_run.txt")
		if err := os.WriteFile(fPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		output := captureStdout(func() {
			err := RenameFiles([]string{fPath}, renameFunc, true, false, true)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		var result RenameResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		if !result.DryRun {
			t.Error("Expected DryRun to be true")
		}
		if result.Summary.Planned != 1 {
			t.Errorf("Expected 1 planned, got %d", result.Summary.Planned)
		}
		if len(result.Operations) != 1 {
			t.Fatalf("Expected 1 operation, got %d", len(result.Operations))
		}
		if result.Operations[0].Status != StatusPlanned {
			t.Errorf("Expected operation status %s, got %s", StatusPlanned, result.Operations[0].Status)
		}
	})

	t.Run("successful execution JSON", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "exec.txt")
		if err := os.WriteFile(fPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		output := captureStdout(func() {
			err := RenameFiles([]string{fPath}, renameFunc, false, false, true)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		var result RenameResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if result.DryRun {
			t.Error("Expected DryRun to be false")
		}
		if result.Summary.Renamed != 1 {
			t.Errorf("Expected 1 renamed, got %d", result.Summary.Renamed)
		}
		if len(result.Operations) != 1 {
			t.Fatalf("Expected 1 operation, got %d", len(result.Operations))
		}
		if result.Operations[0].Status != StatusRenamed {
			t.Errorf("Expected operation status %s, got %s", StatusRenamed, result.Operations[0].Status)
		}
	})

	t.Run("unchanged operation", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "UNCHANGED.txt")
		if err := os.WriteFile(fPath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		output := captureStdout(func() {
			err := RenameFiles([]string{fPath}, renameFunc, false, false, true)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		var result RenameResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v", err)
		}

		if result.Summary.Unchanged != 1 {
			t.Errorf("Expected 1 unchanged, got %d", result.Summary.Unchanged)
		}
		if result.Operations[0].Status != StatusUnchanged {
			t.Errorf("Expected operation status %s, got %s", StatusUnchanged, result.Operations[0].Status)
		}
	})

	t.Run("collision exhaustive batch", func(t *testing.T) {
		fPath1 := filepath.Join(tempDir, "col1.txt")
		fPath2 := filepath.Join(tempDir, "col2.txt")
		fPath3 := filepath.Join(tempDir, "col3.txt") // Will be a collision with col1.txt
		fPath4 := filepath.Join(tempDir, "COL4.txt") // Already uppercase (unchanged)

		if err := os.WriteFile(fPath1, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file 1: %v", err)
		}
		if err := os.WriteFile(fPath2, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file 2: %v", err)
		}
		if err := os.WriteFile(fPath3, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file 3: %v", err)
		}
		if err := os.WriteFile(fPath4, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write test file 4: %v", err)
		}

		// Map col3.txt to col1.txt destination (COL1.TXT)
		badRenameFunc := func(s string) (string, error) {
			if s == "col3" {
				return "COL1", nil
			}
			return strings.ToUpper(s), nil
		}

		output := captureStdout(func() {
			err := RenameFiles([]string{fPath1, fPath2, fPath3, fPath4}, badRenameFunc, false, false, true)
			if err == nil {
				t.Error("Expected error due to collision, got nil")
			}
		})

		var result RenameResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		if result.Summary.Collision != 2 {
			t.Errorf("Expected 2 collisions, got %d", result.Summary.Collision)
		}
		if result.Summary.Skipped != 1 {
			t.Errorf("Expected 1 skipped, got %d", result.Summary.Skipped)
		}
		if result.Summary.Unchanged != 1 {
			t.Errorf("Expected 1 unchanged, got %d", result.Summary.Unchanged)
		}
		if len(result.Operations) != 4 {
			t.Fatalf("Expected 4 operations, got %d", len(result.Operations))
		}

		hasSkipped := false
		hasUnchanged := false
		hasCollision := 0
		for _, op := range result.Operations {
			if op.Status == StatusSkipped {
				hasSkipped = true
			}
			if op.Status == StatusUnchanged {
				hasUnchanged = true
			}
			if op.Status == StatusCollision {
				hasCollision++
			}
		}
		if !hasSkipped {
			t.Error("Expected one skipped operation")
		}
		if !hasUnchanged {
			t.Error("Expected one unchanged operation")
		}
		if hasCollision != 2 {
			t.Errorf("Expected 2 collision operations, got %d", hasCollision)
		}

		// Verify no files were renamed
		if _, err := os.Stat(filepath.Join(tempDir, "COL2.txt")); !os.IsNotExist(err) {
			t.Error("Files were renamed even though batch aborted")
		}
	})

	t.Run("failed rename", func(t *testing.T) {
		fPath := filepath.Join(tempDir, "fail_test.txt")

		output := captureStdout(func() {
			// file doesn't exist, so rename execution should fail, BUT planning preflight might pass if we aren't careful
			// However, in rntocase preflight lstat checks if destination exists, not if source exists.
			// Let's force a failure by providing a rename function that returns an error
			badRenameFunc := func(s string) (string, error) {
				return "", fmt.Errorf("forced error")
			}
			err := RenameFiles([]string{fPath}, badRenameFunc, false, false, true)
			if err == nil {
				t.Error("Expected error due to failed generation, got nil")
			}
		})

		var result RenameResult
		if err := json.Unmarshal([]byte(output), &result); err != nil {
			t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
		}

		if result.Summary.Failed != 1 {
			t.Errorf("Expected 1 failed, got %d", result.Summary.Failed)
		}
		if len(result.Operations) != 1 {
			t.Fatalf("Expected 1 operation, got %d", len(result.Operations))
		}
		if result.Operations[0].Status != StatusFailed {
			t.Errorf("Expected operation status %s, got %s", StatusFailed, result.Operations[0].Status)
		}
	})
}
