package fsys

import (
	"os"
	"path/filepath"
	"sort"
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

func (m *MockFS) stat(name string, depth int) (os.FileInfo, error) {
	if depth > 40 {
		return nil, os.ErrInvalid
	}

	name = filepath.Clean(name)
	if name != "." && name != "/" {
		evalDir, err := m.EvalSymlinks(filepath.Dir(name))
		if err != nil {
			return nil, err
		}
		name = filepath.Join(evalDir, filepath.Base(name))
	}
	fi, ok := m.Files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	// naive evaluate symlinks
	if fi.ModeVal&os.ModeSymlink != 0 {
		target := fi.SymlinkTarget
		if !filepath.IsAbs(target) {
			target = filepath.Join(filepath.Dir(name), target)
		}
		return m.stat(target, depth+1)
	}
	return fi, nil
}

func (m *MockFS) Stat(name string) (os.FileInfo, error) {
	return m.stat(name, 0)
}

func (m *MockFS) Lstat(name string) (os.FileInfo, error) {
	name = filepath.Clean(name)
	if name == "." || name == "/" {
		fi, ok := m.Files[name]
		if ok {
			return fi, nil
		}
		return nil, os.ErrNotExist
	}
	// Resolve ancestor directory
	evalDir, err := m.EvalSymlinks(filepath.Dir(name))
	if err != nil {
		return nil, err
	}
	name = filepath.Join(evalDir, filepath.Base(name))

	if fi, ok := m.Files[name]; ok {
		return fi, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) evalSymlinks(path string, depth int) (string, error) {
	if depth > 40 {
		return "", os.ErrInvalid
	}

	path = filepath.Clean(path)
	if path == "." || path == "/" {
		return path, nil
	}
	dir := filepath.Dir(path)
	evalDir, err := m.evalSymlinks(dir, depth)
	if err != nil {
		return "", err
	}

	path = filepath.Join(evalDir, filepath.Base(path))
	fi, ok := m.Files[path]
	if !ok {
		return "", os.ErrNotExist
	}

	if fi.ModeVal&os.ModeSymlink != 0 {
		target := fi.SymlinkTarget
		if !filepath.IsAbs(target) {
			target = filepath.Join(evalDir, target)
		}
		// don't increase depth here, we only increase depth when following a link
		return m.evalSymlinks(target, depth+1)
	}
	return path, nil
}

func (m *MockFS) EvalSymlinks(path string) (string, error) {
	return m.evalSymlinks(filepath.Clean(path), 0)
}

func (m *MockFS) Rename(oldpath, newpath string) error {
	if oldpath == newpath {
		return nil
	}

	oldpath = filepath.Clean(oldpath)
	newpath = filepath.Clean(newpath)

	evalOldDir, err := m.EvalSymlinks(filepath.Dir(oldpath))
	if err != nil {
		return err
	}
	resolvedOld := filepath.Join(evalOldDir, filepath.Base(oldpath))

	evalNewDir, err := m.EvalSymlinks(filepath.Dir(newpath))
	if err != nil {
		return err
	}
	resolvedNew := filepath.Join(evalNewDir, filepath.Base(newpath))

	fi, ok := m.Files[resolvedOld]
	if !ok {
		return os.ErrNotExist
	}

	if _, destOk := m.Files[resolvedNew]; destOk {
		return os.ErrExist
	}

	destDir := evalNewDir
	if destDir != "." && destDir != "/" {
		destDirFi, destDirErr := m.stat(destDir, 0)
		if destDirErr != nil {
			return destDirErr
		}
		if !destDirFi.IsDir() {
			return os.ErrInvalid
		}
	}

	delete(m.Files, resolvedOld)
	fi.NameStr = filepath.Base(resolvedNew)
	m.Files[resolvedNew] = fi
	return nil
}

func (m *MockFS) SameFile(fi1, fi2 os.FileInfo) bool {
	// MockFS identity: compare names and modes for simplicity, or just pointers if we return same pointers.
	// Our MockFS Stat/Lstat return pointers to the internal MockFileInfo, so pointer equality works.
	return fi1 == fi2
}

func (m *MockFS) ReadDir(name string) ([]os.DirEntry, error) {
	name = filepath.Clean(name)
	resolvedName, err := m.EvalSymlinks(name)
	if err != nil {
		return nil, err
	}
	name = resolvedName
	fi, err := m.Lstat(name) // it is resolved now, so Lstat is fine, and we check if it is dir
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, os.ErrInvalid
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
