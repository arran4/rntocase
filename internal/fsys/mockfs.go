package fsys

import (
	"golang.org/x/tools/txtar"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MockFS struct {
	Files map[string]*MockFileInfo
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string]*MockFileInfo)}
}

type MockFileInfo struct {
	NameStr       string
	ModeVal       os.FileMode
	IsDirVal      bool
	SymlinkTarget string
}

func (m *MockFileInfo) Name() string               { return m.NameStr }
func (m *MockFileInfo) Size() int64                { return 0 }
func (m *MockFileInfo) Mode() os.FileMode          { return m.ModeVal }
func (m *MockFileInfo) ModTime() time.Time         { return time.Time{} }
func (m *MockFileInfo) IsDir() bool                { return m.IsDirVal }
func (m *MockFileInfo) Sys() any                   { return nil }
func (m *MockFileInfo) Info() (os.FileInfo, error) { return m, nil }
func (m *MockFileInfo) Type() os.FileMode          { return m.ModeVal.Type() }

func (m *MockFS) AddFile(name string) {
	name = filepath.Clean(name)
	m.Files[name] = &MockFileInfo{NameStr: filepath.Base(name), ModeVal: 0644, IsDirVal: false}
	m.ensureDirs(filepath.Dir(name))
}

func (m *MockFS) AddDir(name string) {
	name = filepath.Clean(name)
	m.Files[name] = &MockFileInfo{NameStr: filepath.Base(name), ModeVal: 0755 | os.ModeDir, IsDirVal: true}
	m.ensureDirs(filepath.Dir(name))
}

func (m *MockFS) AddSymlink(name, target string) {
	name = filepath.Clean(name)
	m.Files[name] = &MockFileInfo{NameStr: filepath.Base(name), ModeVal: 0777 | os.ModeSymlink, IsDirVal: false, SymlinkTarget: target}
	m.ensureDirs(filepath.Dir(name))
}

func (m *MockFS) ensureDirs(dir string) {
	if dir == "." || dir == "/" || dir == "" {
		return
	}
	if _, ok := m.Files[dir]; !ok {
		m.Files[dir] = &MockFileInfo{NameStr: filepath.Base(dir), ModeVal: 0755 | os.ModeDir, IsDirVal: true}
		m.ensureDirs(filepath.Dir(dir))
	}
}

func (m *MockFS) Stat(name string) (os.FileInfo, error) {
	name = filepath.Clean(name)
	fi, ok := m.Files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	// naive evaluate symlinks
	if fi.ModeVal&os.ModeSymlink != 0 {
		return m.Stat(filepath.Join(filepath.Dir(name), fi.SymlinkTarget))
	}
	return fi, nil
}

func (m *MockFS) Lstat(name string) (os.FileInfo, error) {
	name = filepath.Clean(name)
	if fi, ok := m.Files[name]; ok {
		return fi, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) EvalSymlinks(path string) (string, error) {
	path = filepath.Clean(path)
	fi, ok := m.Files[path]
	if !ok {
		return "", os.ErrNotExist
	}
	if fi.ModeVal&os.ModeSymlink != 0 {
		return filepath.Clean(filepath.Join(filepath.Dir(path), fi.SymlinkTarget)), nil
	}
	return path, nil
}

func (m *MockFS) Rename(oldpath, newpath string) error {
	if oldpath == newpath {
		return nil
	}

	oldpath = filepath.Clean(oldpath)
	newpath = filepath.Clean(newpath)

	fi, ok := m.Files[oldpath]
	if !ok {
		return os.ErrNotExist
	}

	if _, destOk := m.Files[newpath]; destOk {
		return os.ErrExist
	}

	delete(m.Files, oldpath)
	fi.NameStr = filepath.Base(newpath)
	m.Files[newpath] = fi
	return nil
}

func (m *MockFS) SameFile(fi1, fi2 os.FileInfo) bool {
	// MockFS identity: compare names and modes for simplicity, or just pointers if we return same pointers.
	// Our MockFS Stat/Lstat return pointers to the internal MockFileInfo, so pointer equality works.
	return fi1 == fi2
}

func (m *MockFS) ReadDir(name string) ([]os.DirEntry, error) {
	name = filepath.Clean(name)
	_, err := m.Lstat(name)
	if err != nil {
		return nil, err
	}

	var entries []os.DirEntry
	for k, fi := range m.Files {
		if filepath.Dir(k) == name && k != name {
			entries = append(entries, fi)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return entries, nil
}

func ArchiveToMockFS(ar *txtar.Archive) *MockFS {
	out := NewMockFS()
	for _, f := range ar.Files {
		name := filepath.Clean(strings.TrimPrefix(f.Name, "/"))
		if name == "." {
			continue
		}

		// simple way to encode directories: name ends with /
		if strings.HasSuffix(f.Name, "/") {
			out.AddDir(name)
		} else {
			out.AddFile(name)
			// Can be extended to support symlinks by checking comment or content
			if strings.HasPrefix(string(f.Data), "SYMLINK:") {
				target := strings.TrimSpace(strings.TrimPrefix(string(f.Data), "SYMLINK:"))
				out.AddSymlink(name, target)
			}
		}
	}
	return out
}
