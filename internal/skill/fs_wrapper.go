package skill

import "os"

// FSOps defines the filesystem operations needed for safe replacement.
// This allows us to use an in-memory fake for most unit tests.
type FSOps interface {
	MkdirAll(path string, perm os.FileMode) error
	MkdirTemp(dir, pattern string) (string, error)
	Remove(name string) error
	RemoveAll(path string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (os.FileInfo, error)
}

// OSFSOps provides the real OS filesystem operations.
type OSFSOps struct{}

func (OSFSOps) MkdirAll(path string, perm os.FileMode) error  { return os.MkdirAll(path, perm) }
func (OSFSOps) MkdirTemp(dir, pattern string) (string, error) { return os.MkdirTemp(dir, pattern) }
func (OSFSOps) Remove(name string) error                      { return os.Remove(name) }
func (OSFSOps) RemoveAll(path string) error                   { return os.RemoveAll(path) }
func (OSFSOps) Rename(oldpath, newpath string) error          { return os.Rename(oldpath, newpath) }
func (OSFSOps) Stat(name string) (os.FileInfo, error)         { return os.Stat(name) }
