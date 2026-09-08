package skill

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type FakeFSOps struct {
	mu           sync.Mutex
	Dirs         map[string]bool
	Files        map[string]bool
	RenameFailFS map[string]error
	StatFailFS   map[string]error
	MkdirTempFn  func(dir, pattern string) (string, error)
	tempCounter  int
}

func NewFakeFSOps() *FakeFSOps {
	return &FakeFSOps{
		Dirs:         make(map[string]bool),
		Files:        make(map[string]bool),
		RenameFailFS: make(map[string]error),
		StatFailFS:   make(map[string]error),
	}
}

func (f *FakeFSOps) MkdirAll(path string, perm os.FileMode) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.Dirs[path] = true
	return nil
}

func (f *FakeFSOps) MkdirTemp(dir, pattern string) (string, error) {
	if f.MkdirTempFn != nil {
		return f.MkdirTempFn(dir, pattern)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tempCounter++
	name := fmt.Sprintf("%s/%s%d", dir, strings.ReplaceAll(pattern, "*", ""), f.tempCounter)
	f.Dirs[name] = true
	return name, nil
}

func (f *FakeFSOps) Remove(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.Dirs, name)
	delete(f.Files, name)
	return nil
}

func (f *FakeFSOps) RemoveAll(path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.Dirs, path)
	for k := range f.Dirs {
		if strings.HasPrefix(k, path+"/") {
			delete(f.Dirs, k)
		}
	}
	for k := range f.Files {
		if strings.HasPrefix(k, path+"/") {
			delete(f.Files, k)
		}
	}
	return nil
}

func (f *FakeFSOps) Rename(oldpath, newpath string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err, ok := f.RenameFailFS[oldpath+"->"+newpath]; ok {
		return err
	}

	// Simple rename logic
	if f.Dirs[oldpath] {
		delete(f.Dirs, oldpath)
		f.Dirs[newpath] = true
		// Rename children
		for k := range f.Dirs {
			if strings.HasPrefix(k, oldpath+"/") {
				f.Dirs[strings.Replace(k, oldpath, newpath, 1)] = true
				delete(f.Dirs, k)
			}
		}
		for k := range f.Files {
			if strings.HasPrefix(k, oldpath+"/") {
				f.Files[strings.Replace(k, oldpath, newpath, 1)] = true
				delete(f.Files, k)
			}
		}
	} else if f.Files[oldpath] {
		delete(f.Files, oldpath)
		f.Files[newpath] = true
	} else {
		return os.ErrNotExist
	}

	return nil
}

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return 0 }
func (f fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f fakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fakeFileInfo) IsDir() bool        { return f.isDir }
func (f fakeFileInfo) Sys() any           { return nil }

func (f *FakeFSOps) Stat(name string) (os.FileInfo, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err, ok := f.StatFailFS[name]; ok {
		return nil, err
	}

	if f.Dirs[name] {
		return fakeFileInfo{name: name, isDir: true}, nil
	}
	if f.Files[name] {
		return fakeFileInfo{name: name, isDir: false}, nil
	}
	return nil, os.ErrNotExist
}
