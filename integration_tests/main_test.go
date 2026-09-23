package integration_tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var sharedBinPath string

func TestMain(m *testing.M) {
	tempDir, err := os.MkdirTemp("", "rntocase_integration_tests_*")
	if err != nil {
		os.Exit(1)
	}

	sharedBinPath = filepath.Join(tempDir, "rntocase")
	cmd := exec.Command("go", "build", "-o", sharedBinPath, "github.com/arran4/rntocase/cmd/rntocase")
	cmd.Dir = "../"
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(tempDir)
		os.Exit(1)
	}

	code := m.Run()

	_ = os.RemoveAll(tempDir)
	os.Exit(code)
}
