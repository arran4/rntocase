package skill

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GitHubCommit represents the basic info of a GitHub commit.
type GitHubCommit struct {
	Sha string `json:"sha"`
}

// DownloadGitHubRepository downloads a repository tarball and returns the temp file path and commit SHA.
func DownloadGitHubRepository(ownerRepo, ref string) (string, string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	// 1. Get latest commit SHA to track revision
	apiURL := fmt.Sprintf("%s/repos/%s/commits/%s", GitHubAPIURL, ownerRepo, ref)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	// Use github token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch repository metadata: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to get repository metadata: HTTP %d", resp.StatusCode)
	}

	var commit GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commit); err != nil {
		return "", "", fmt.Errorf("failed to decode repository metadata: %w", err)
	}

	// 2. Download tarball
	tarballURL := fmt.Sprintf("%s/repos/%s/tarball/%s", GitHubAPIURL, ownerRepo, commit.Sha)
	reqTar, err := http.NewRequest("GET", tarballURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create tarball request: %w", err)
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		reqTar.Header.Set("Authorization", "token "+token)
	}

	respTar, err := HTTPClient.Do(reqTar)
	if err != nil {
		return "", "", fmt.Errorf("failed to download tarball: %w", err)
	}
	defer func() { _ = respTar.Body.Close() }()

	if respTar.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to download tarball: HTTP %d", respTar.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "rntocase-skill-*.tar.gz")
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() { _ = tmpFile.Close() }()

	if _, err := io.Copy(tmpFile, respTar.Body); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", "", fmt.Errorf("failed to write tarball: %w", err)
	}

	return tmpFile.Name(), commit.Sha, nil
}

// ExtractTarGz extracts the given tarball to the destination directory.
func ExtractTarGz(tarballPath, destDir, pathWithin string) error {
	file, err := os.Open(tarballPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer func() { _ = file.Close() }()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)

	// GitHub tarballs have a top-level directory like owner-repo-sha/
	var topLevelDir string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		if topLevelDir == "" {
			parts := strings.SplitN(header.Name, "/", 2)
			if len(parts) > 0 {
				topLevelDir = parts[0]
			}
		}

		// Normalize file path inside the tar
		relativePath := strings.TrimPrefix(header.Name, topLevelDir+"/")

		// If pathWithin is provided, only extract files under that path
		if pathWithin != "" {
			if !strings.HasPrefix(relativePath, pathWithin+"/") && relativePath != pathWithin {
				continue
			}
			// Strip pathWithin prefix for the destination
			relativePath = strings.TrimPrefix(relativePath, pathWithin)
			relativePath = strings.TrimPrefix(relativePath, "/")
			if relativePath == "" {
				continue // Skip the directory itself if we only want its contents
			}
		}

		// Ensure no empty string or absolute paths
		if relativePath == "" {
			continue
		}

		if err := extractAndProtect(tr, header, destDir, relativePath); err != nil {
			return err
		}
	}

	return nil
}

// extractAndProtect handles the extraction of a single tar entry with path traversal protection.
func extractAndProtect(tr *tar.Reader, header *tar.Header, destDir string, relativePath string) error {
	// SECURITY: Path traversal protection
	targetPath := filepath.Join(destDir, relativePath)
	if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path in tarball (path traversal detected): %s", header.Name)
	}

	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	case tar.TypeReg:
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
		outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer func() { _ = outFile.Close() }()

		if _, err := io.Copy(outFile, tr); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	case tar.TypeSymlink:
		// SECURITY: Prevent malicious symlinks outside the destination
		symlinkTarget := header.Linkname
		absSymlinkTarget := targetPath
		if filepath.IsAbs(symlinkTarget) {
			return fmt.Errorf("absolute symlinks are not allowed: %s", header.Name)
		}

		resolvedSymlink := filepath.Join(filepath.Dir(absSymlinkTarget), symlinkTarget)
		if !strings.HasPrefix(filepath.Clean(resolvedSymlink), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("symlink points outside destination directory: %s", header.Name)
		}

		if err := os.Symlink(symlinkTarget, targetPath); err != nil {
			return fmt.Errorf("failed to create symlink: %w", err)
		}
	}

	return nil
}

// CopyLocalDirectory copies a local directory to the destination safely.
func CopyLocalDirectory(src, dest string) error {
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)

	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil // Skip the root directory itself
		}

		destPath := filepath.Join(dest, relPath)

		// Basic path traversal protection for local copies too
		if !strings.HasPrefix(destPath, dest+string(os.PathSeparator)) {
			return fmt.Errorf("path traversal detected during local copy: %s", destPath)
		}

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		if info.Mode()&os.ModeSymlink != 0 {
			symlinkTarget, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(symlinkTarget, destPath)
		}

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()

		destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer func() { _ = destFile.Close() }()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}

		return nil
	})
}

// ReplaceSafely creates a staging directory, populates it using the provided function,
// and if successful, atomically (or via safe backup-and-rename) swaps it into destDir.
// If populateFunc fails, destDir is untouched.
// If the final rename fails, destDir is rolled back to its previous state if possible.
func ReplaceSafely(destDir string, populateFunc func(stagingDir string) error) error {
	return replaceSafelyWithFS(destDir, populateFunc, OSFSOps{})
}

func replaceSafelyWithFS(destDir string, populateFunc func(stagingDir string) error, fsOps FSOps) error {
	parentDir := filepath.Dir(destDir)
	if err := fsOps.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// 1. Create a staging directory
	stagingDir, err := fsOps.MkdirTemp(parentDir, ".staging-skill-*")
	if err != nil {
		return fmt.Errorf("failed to create staging directory: %w", err)
	}
	defer func() {
		_ = fsOps.RemoveAll(stagingDir)
	}()

	// 2. Populate the staging directory
	if err := populateFunc(stagingDir); err != nil {
		return err // Cleanup handled by defer
	}

	// 3. Swap the directories safely
	backupDir := ""
	if _, err := fsOps.Stat(destDir); err == nil {
		// Existing directory found, create a backup
		backupDir, err = fsOps.MkdirTemp(parentDir, ".backup-skill-*")
		if err != nil {
			return fmt.Errorf("failed to create backup directory: %w", err)
		}
		// MkdirTemp creates it, so we need to remove it so Rename can use the path, or just create the name.
		// os.Rename requires the destination to be empty or not exist on most platforms.
		_ = fsOps.Remove(backupDir)

		if err := fsOps.Rename(destDir, backupDir); err != nil {
			return fmt.Errorf("failed to backup existing installation: %w", err)
		}
	}

	// Move staging to destination
	if err := fsOps.Rename(stagingDir, destDir); err != nil {
		// Rollback on failure
		if backupDir != "" {
			if restoreErr := fsOps.Rename(backupDir, destDir); restoreErr != nil {
				return fmt.Errorf("failed to commit new installation (%v) and rollback failed: %v", err, restoreErr)
			}
		}
		return fmt.Errorf("failed to commit new installation: %w", err)
	}

	// Cleanup backup on success
	if backupDir != "" {
		_ = fsOps.RemoveAll(backupDir)
	}

	return nil
}
