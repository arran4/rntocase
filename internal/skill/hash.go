package skill

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// ComputeDirectoryDigest computes a combined SHA256 digest of all files in the directory
// (ignoring our metadata file) to detect local modifications.
func ComputeDirectoryDigest(dir string) (string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if relPath == MetadataFileName {
			return nil
		}

		files = append(files, relPath)
		return nil
	})
	if err != nil {
		return "", err
	}

	sort.Strings(files)

	h := sha256.New()
	for _, relPath := range files {
		path := filepath.Join(dir, relPath)
		h.Write([]byte(relPath))
		h.Write([]byte{0}) // Delimiter to prevent collision

		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(h, f); err != nil {
			_ = f.Close()
			return "", err
		}
		_ = f.Close()
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
