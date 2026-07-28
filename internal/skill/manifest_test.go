package skill

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoadMetadata(t *testing.T) {
	dir, err := os.MkdirTemp("", "skill-meta-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	now := time.Now().Truncate(time.Second) // JSON serialization truncates to some precision

	meta := &Metadata{
		Name:           "test-skill",
		OriginalSource: "owner/repo",
		OwnerRepo:      "owner/repo",
		SourceRevision: "abcdef123456",
		InstallTime:    now,
		InstallerApp:   "test-app",
	}

	err = SaveMetadata(dir, meta)
	assert.NoError(t, err)

	assert.FileExists(t, filepath.Join(dir, MetadataFileName))

	loaded, err := LoadMetadata(dir)
	assert.NoError(t, err)

	assert.Equal(t, meta.Name, loaded.Name)
	assert.Equal(t, meta.OriginalSource, loaded.OriginalSource)
	assert.Equal(t, meta.OwnerRepo, loaded.OwnerRepo)
	assert.Equal(t, meta.SourceRevision, loaded.SourceRevision)
	assert.Equal(t, meta.InstallerApp, loaded.InstallerApp)
	assert.True(t, now.Equal(loaded.InstallTime), "time mismatch")
}

func TestLoadMetadata_NotFound(t *testing.T) {
	dir, err := os.MkdirTemp("", "skill-meta-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(dir) }()

	_, err = LoadMetadata(dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "metadata not found")
}
