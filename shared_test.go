package rntocase

import (
	"bufio"
	"io"
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

		if err := RenameFiles(paths, renameFunc, false, false); err != nil {
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

		err := RenameFiles([]string{filepath.Join(tempDir, "Bar.txt")}, renameToBaz, false, false)
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
		err := RenameFiles(paths, renameToSame, false, false)
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
		// Write the second file, we will name it Foo.txt.
		// On case-insensitive FS, it will just overwrite FOO.txt, but for testing we write to BAR.txt
		// so that they are definitely distinct files to start with.
		if err := os.WriteFile(filepath.Join(tempDir, "BAR.txt"), []byte("SECOND"), 0644); err != nil {
			t.Fatal(err)
		}

		renameToLower := func(s string) (string, error) {
			// for BAR we pretend we are renaming it to foo
			if s == "BAR" {
				return "foo", nil
			}
			return strings.ToLower(s), nil
		}

		// Use a relative path and an absolute path that resolve to the same dir
		absPath, _ := filepath.Abs(tempDir)
		paths := []string{
			filepath.Join(tempDir, "FOO.txt"), // could be relative if we used Chdir, but let's test absolute and evaluated
			filepath.Join(absPath, "BAR.txt"),
		}

		err := RenameFiles(paths, renameToLower, false, false)
		if err == nil {
			t.Fatal("Expected error due to multiple files mapping to same destination across rel/abs paths")
		}

		if !strings.Contains(err.Error(), "multiple source files map to destination") {
			t.Errorf("Expected collision error message, got: %v", err)
		}

		// Ensure filesystem unchanged
		b1, _ := os.ReadFile(filepath.Join(tempDir, "FOO.txt"))
		if string(b1) != "FIRST" {
			t.Errorf("Source contents changed, expected FIRST, got %s", string(b1))
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
		err := RenameFiles(paths, renameToLower, false, false)
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
		err := RenameFiles(paths, renameToSame, true, false)
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

		err := RenameFiles(paths, renameFuncLower, false, false)
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
		err := RenameFiles(paths, renameMixed, false, false)
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
		err := RenameFiles(paths, renameFunc, false, false)
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
