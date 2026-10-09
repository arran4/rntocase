package fsys

import (
	"os"
	"strings"
	"testing"
)

func TestMockFS_EvalSymlinks_Cycles(t *testing.T) {
	m := NewMockFS()

	t.Run("nested-parent cycle", func(t *testing.T) {
		m.AddSymlink("/a/link", "/a/link/file")
		_, err := m.EvalSymlinks("/a/link")
		if err == nil {
			t.Fatal("Expected error on cyclic symlink, got nil")
		}
		if err != os.ErrInvalid {
			t.Fatalf("Expected os.ErrInvalid, got %v", err)
		}
	})

	t.Run("straightforward direct cycle", func(t *testing.T) {
		m.AddSymlink("/cycle", "/cycle")
		_, err := m.EvalSymlinks("/cycle")
		if err != os.ErrInvalid {
			t.Fatalf("Expected os.ErrInvalid, got %v", err)
		}
	})

	t.Run("multi-link cycle", func(t *testing.T) {
		m.AddSymlink("/cycle1", "/cycle2")
		m.AddSymlink("/cycle2", "/cycle1")
		_, err := m.EvalSymlinks("/cycle1")
		if err != os.ErrInvalid {
			t.Fatalf("Expected os.ErrInvalid, got %v", err)
		}
	})

	t.Run("deep directory structure without cycles", func(t *testing.T) {
		// Create a deep directory structure
		path := "/root"
		for i := 0; i < 50; i++ {
			path += "/dir"
		}
		m.AddDir(path)

		resolved, err := m.EvalSymlinks(path)
		if err != nil {
			t.Fatalf("Expected deep directory structure to resolve successfully, got %v", err)
		}
		if resolved != path {
			t.Fatalf("Expected path %s, got %s", path, resolved)
		}
	})
}

func TestMockFS_Aliases(t *testing.T) {
	m := NewMockFS()

	// Create an absolute ancestor alias and a relative ancestor alias
	m.AddDir("/target/sub")
	m.AddFile("/target/sub/file.txt")

	m.AddSymlink("/alias_abs", "/target")
	m.AddSymlink("/alias_rel", "target") // relative to /

	t.Run("ReadDir through absolute alias", func(t *testing.T) {
		entries, err := m.ReadDir("/alias_abs/sub")
		if err != nil {
			t.Fatalf("Expected success reading aliased directory, got %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "file.txt" {
			t.Fatalf("Expected [file.txt], got %v", entries)
		}
	})

	t.Run("ReadDir through relative alias", func(t *testing.T) {
		entries, err := m.ReadDir("/alias_rel/sub")
		if err != nil {
			t.Fatalf("Expected success reading aliased directory, got %v", err)
		}
		if len(entries) != 1 || entries[0].Name() != "file.txt" {
			t.Fatalf("Expected [file.txt], got %v", entries)
		}
	})

	t.Run("Stat through alias", func(t *testing.T) {
		fi, err := m.Stat("/alias_abs/sub/file.txt")
		if err != nil {
			t.Fatalf("Expected success statting aliased file, got %v", err)
		}
		if fi.Name() != "file.txt" {
			t.Fatalf("Expected file.txt, got %v", fi.Name())
		}
	})

	t.Run("Lstat through alias", func(t *testing.T) {
		fi, err := m.Lstat("/alias_rel/sub/file.txt")
		if err != nil {
			t.Fatalf("Expected success Lstatting aliased file, got %v", err)
		}
		if fi.Name() != "file.txt" {
			t.Fatalf("Expected file.txt, got %v", fi.Name())
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("Expected regular file, got symlink")
		}
	})

	t.Run("Rename and collision behavior through alias", func(t *testing.T) {
		m.AddFile("/other.txt")
		err := m.Rename("/other.txt", "/alias_abs/sub/file.txt") // should collide
		if err == nil {
			t.Fatalf("Expected collision renaming to existing aliased file, got nil")
		}

		err = m.Rename("/other.txt", "/alias_abs/sub/newfile.txt") // should succeed
		if err != nil {
			t.Fatalf("Expected success renaming to aliased directory, got %v", err)
		}

		// Verify original file is gone and new file exists through alias
		_, err = m.Lstat("/other.txt")
		if err == nil {
			t.Fatalf("Expected source file to be removed after rename")
		}
		_, err = m.Lstat("/target/sub/newfile.txt")
		if err != nil {
			t.Fatalf("Expected new file to exist in target directory")
		}
	})

	t.Run("dangling links", func(t *testing.T) {
		m.AddSymlink("/dangling", "/nonexistent")

		// Lstat should succeed
		fi, err := m.Lstat("/dangling")
		if err != nil {
			t.Fatalf("Expected Lstat to succeed on dangling link, got %v", err)
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("Expected symlink mode, got %v", fi.Mode())
		}

		// Stat should fail
		_, err = m.Stat("/dangling")
		if err == nil {
			t.Fatalf("Expected Stat to fail on dangling link, got nil")
		}

		// EvalSymlinks should fail
		_, err = m.EvalSymlinks("/dangling")
		if err == nil {
			t.Fatalf("Expected EvalSymlinks to fail on dangling link, got nil")
		}
	})

	t.Run("ReadDir on non-directory should error", func(t *testing.T) {
		m.AddFile("/file.txt")
		_, err := m.ReadDir("/file.txt")
		if err == nil || !strings.Contains(err.Error(), "invalid") {
			t.Fatalf("Expected invalid directory error on ReadDir(/file.txt), got %v", err)
		}
	})
}
