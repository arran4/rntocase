package skill

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTestTarball(t *testing.T, files map[string]string) string {
	t.Helper()
	f, err := os.CreateTemp("", "test-*.tar.gz")
	assert.NoError(t, err)

	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	for name, content := range files {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	assert.NoError(t, tw.Close())
	assert.NoError(t, gw.Close())
	assert.NoError(t, f.Close())

	return f.Name()
}

func TestExtractTarGz_PathTraversal(t *testing.T) {
	destDir, err := os.MkdirTemp("", "dest-dir-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(destDir) }()

	files := map[string]string{
		"repo-sha/valid.txt":      "valid content",
		"repo-sha/../../evil.txt": "evil content",
	}

	tarPath := createTestTarball(t, files)
	defer func() { _ = os.Remove(tarPath) }()

	err = ExtractTarGz(tarPath, destDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal detected")
}

func TestExtractTarGz_Success(t *testing.T) {
	destDir, err := os.MkdirTemp("", "dest-dir-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(destDir) }()

	files := map[string]string{
		"repo-sha/SKILL.md":      "# My Skill\n",
		"repo-sha/lib/helper.py": "print('hello')",
	}

	tarPath := createTestTarball(t, files)
	defer func() { _ = os.Remove(tarPath) }()

	err = ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "# My Skill\n", string(content))
}

func TestReplaceSafely_Success(t *testing.T) {
	parentDir, err := os.MkdirTemp("", "parent-dir-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(parentDir) }()

	destDir := filepath.Join(parentDir, "my-skill")

	// Initially does not exist
	err = ReplaceSafely(destDir, func(stagingDir string) error {
		return os.WriteFile(filepath.Join(stagingDir, "SKILL.md"), []byte("v1"), 0644)
	})
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "v1", string(content))

	// Replace existing
	err = ReplaceSafely(destDir, func(stagingDir string) error {
		return os.WriteFile(filepath.Join(stagingDir, "SKILL.md"), []byte("v2"), 0644)
	})
	assert.NoError(t, err)

	content, err = os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "v2", string(content))
}

func TestReplaceSafely_RollbackOnFailure(t *testing.T) {
	parentDir, err := os.MkdirTemp("", "parent-dir-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(parentDir) }()

	destDir := filepath.Join(parentDir, "my-skill")
	err = os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("working-v1"), 0644)
	assert.NoError(t, err)

	// Simulate a failure during populate (e.g., validation failed)
	err = ReplaceSafely(destDir, func(stagingDir string) error {
		return os.ErrPermission // Some simulated error
	})
	assert.Error(t, err)

	// Original destination should remain completely untouched
	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "working-v1", string(content))
}

func TestReplaceSafely_RemovesStaleFiles(t *testing.T) {
	parentDir, err := os.MkdirTemp("", "parent-dir-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(parentDir) }()

	destDir := filepath.Join(parentDir, "my-skill")
	err = os.MkdirAll(destDir, 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(destDir, "SKILL.md"), []byte("v1"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(destDir, "stale.py"), []byte("print('stale')"), 0644)
	assert.NoError(t, err)

	// Replace existing but new version doesn't have stale.py
	err = ReplaceSafely(destDir, func(stagingDir string) error {
		return os.WriteFile(filepath.Join(stagingDir, "SKILL.md"), []byte("v2"), 0644)
	})
	assert.NoError(t, err)

	// Verify old files are gone
	_, err = os.Stat(filepath.Join(destDir, "stale.py"))
	assert.True(t, os.IsNotExist(err))

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "v2", string(content))
}
