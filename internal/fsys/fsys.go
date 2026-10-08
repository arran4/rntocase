package fsys

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FS interface {
	Lstat(name string) (os.FileInfo, error)
	ReadDir(name string) ([]os.DirEntry, error)
	Stat(name string) (os.FileInfo, error)
	EvalSymlinks(path string) (string, error)
}

type WritableFS interface {
	FS
	Rename(oldpath, newpath string) error
	SameFile(fi1, fi2 os.FileInfo) bool
}

type OSFS struct{}

func (OSFS) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (OSFS) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSFS) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFS) EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

func (OSFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (OSFS) SameFile(fi1, fi2 os.FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

// WalkDir provides a testable filesystem traversal using our FS interface.
func WalkDir(fsys FS, root string, fn fs.WalkDirFunc) error {
	info, err := fsys.Lstat(root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = walkDir(fsys, root, fs.FileInfoToDirEntry(info), fn)
	}
	if err == fs.SkipDir || err == filepath.SkipDir {
		return nil
	}
	return err
}

func walkDir(fsys FS, name string, d fs.DirEntry, fn fs.WalkDirFunc) error {
	if err := fn(name, d, nil); err != nil || !d.IsDir() {
		if err == fs.SkipDir && d.IsDir() {
			err = nil
		}
		return err
	}
	dirs, err := fsys.ReadDir(name)
	if err != nil {
		err = fn(name, d, err)
		if err != nil {
			if err == fs.SkipDir && d.IsDir() {
				err = nil
			}
			return err
		}
	}
	for _, d1 := range dirs {
		name1 := filepath.Join(name, d1.Name())
		if err := walkDir(fsys, name1, d1, fn); err != nil {
			if err == fs.SkipDir {
				continue
			}
			return err
		}
	}
	return nil
}
