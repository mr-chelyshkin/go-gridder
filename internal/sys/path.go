package sys

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"syscall"
)

type path struct {
	dir string

	perm os.FileMode
}

// WithPerm option: add custom perms for directory.
func WithPerm(perm os.FileMode) PathOpt {
	return func(p *path) {
		p.perm = perm
	}
}

type PathOpt func(p *path)

// DirCreate directory.
func DirCreate(dir string, opts ...PathOpt) error {
	p := &path{dir: dir}

	p.perm = 0644
	for _, opt := range opts {
		opt(p)
	}
	if _, err := os.Stat(p.dir); errors.Is(err, os.ErrNotExist) {
		return os.MkdirAll(p.dir, p.perm)
	}
	return nil
}

// DirIsWritable check dir for write access, return bool answer and error description.
func DirIsWritable(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return false, err
	}

	if !info.IsDir() {
		return false, fmt.Errorf("%s is not directory", dir)
	}
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		u, _ := user.Current()
		return false, fmt.Errorf("%s write permission bit is not set for user: %s", dir, u.Name)
	}

	var stat syscall.Stat_t
	if err = syscall.Stat(dir, &stat); err != nil {
		return false, fmt.Errorf("unable to get stat for %s", dir)
	}
	return true, nil
}

// PathIsFile return bool.
func PathIsFile(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

// PathIsDir return bool.
func PathIsDir(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}
	if !fi.IsDir() {
		return false
	}
	return true
}
