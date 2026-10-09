package fstestutil

import (
	"github.com/arran4/rntocase/internal/fsys"
	"golang.org/x/tools/txtar"
	"path/filepath"
	"strings"
)

func ArchiveToMockFS(ar *txtar.Archive) *fsys.MockFS {
	out := fsys.NewMockFS()
	for _, f := range ar.Files {
		name := filepath.Clean(strings.TrimPrefix(f.Name, "/"))
		if name == "." {
			continue
		}

		if strings.HasPrefix(f.Name, "fs/") {
			name := "/" + strings.TrimPrefix(f.Name, "fs/")
			if strings.HasSuffix(name, "/") {
				if strings.HasPrefix(string(f.Data), "SYMLINK: ") {
					out.AddSymlink(name[:len(name)-1], strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK: ")))
				} else {
					out.AddDir(name)
				}
			} else {
				if strings.HasPrefix(string(f.Data), "SYMLINK: ") {
					out.AddSymlink(name, strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK: ")))
				} else {
					out.AddFile(name)
				}
			}
		} else {
			// This matches original behavior
			if strings.HasSuffix(f.Name, "/") {
				out.AddDir(name)
			} else {
				out.AddFile(name)
				if strings.HasPrefix(string(f.Data), "SYMLINK:") {
					target := strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK:"))
					out.AddSymlink(name, target)
				}
			}
		}
	}
	return out
}
