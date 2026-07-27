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
	defer os.RemoveAll(destDir)

	files := map[string]string{
		"repo-sha/valid.txt":               "valid content",
		"repo-sha/../../evil.txt":          "evil content",
	}

	tarPath := createTestTarball(t, files)
	defer os.Remove(tarPath)

	err = ExtractTarGz(tarPath, destDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal detected")
}

func TestExtractTarGz_Success(t *testing.T) {
	destDir, err := os.MkdirTemp("", "dest-dir-*")
	assert.NoError(t, err)
	defer os.RemoveAll(destDir)

	files := map[string]string{
		"repo-sha/SKILL.md": "# My Skill\n",
		"repo-sha/lib/helper.py": "print('hello')",
	}

	tarPath := createTestTarball(t, files)
	defer os.Remove(tarPath)

	err = ExtractTarGz(tarPath, destDir, "")
	assert.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(destDir, "SKILL.md"))
	assert.NoError(t, err)
	assert.Equal(t, "# My Skill\n", string(content))
}
