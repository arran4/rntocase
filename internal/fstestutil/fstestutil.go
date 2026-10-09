package fstestutil

import (
	"github.com/arran4/rntocase/internal/fsys"
	"golang.org/x/tools/txtar"
	"strings"
)

func ArchiveToMockFS(ar *txtar.Archive) *fsys.MockFS {
	out := fsys.NewMockFS()
	for _, f := range ar.Files {
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
		}
	}
	return out
}
