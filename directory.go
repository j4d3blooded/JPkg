package jpkg

import (
	"io"
	"io/fs"
	"time"
)

type jpkgDirOpenerInfo struct {
	ChildPaths []string
	name       string
	path       string
}

type JPkgDir struct {
	pkg    *JPkg
	name   string
	path   string
	dirIdx int
}

func (j *JPkgDir) Close() error {
	return nil
}

func (j *JPkgDir) Read([]byte) (int, error) {
	return 0, fs.ErrInvalid
}

func (j *JPkgDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if n > 0 {
		items, err := j.pkg.ReadDir(j.path)
		if err != nil {
			return nil, err
		}
		if len(items) <= j.dirIdx {
			return nil, io.EOF
		}
		items = items[j.dirIdx:]
		if len(items) == 0 {
			return nil, io.EOF
		}
		if n >= len(items) {
			j.dirIdx = len(items)
			return items, nil
		}
		j.dirIdx += n
		return items[:n], nil
	}

	j.dirIdx = 0
	return j.pkg.ReadDir(j.path)
}

func (j *JPkgDir) Stat() (fs.FileInfo, error) {
	return j, nil
}

// IsDir implements fs.FileInfo.
func (j *JPkgDir) IsDir() bool {
	return true
}

// ModTime implements fs.FileInfo.
func (j *JPkgDir) ModTime() time.Time {
	return j.pkg.packagedAt
}

// Mode implements fs.FileInfo.
func (j *JPkgDir) Mode() fs.FileMode {
	return fs.ModeDir
}

// Name implements fs.FileInfo.
func (j *JPkgDir) Name() string {
	return j.name
}

// Size implements fs.FileInfo.
func (j *JPkgDir) Size() int64 {
	return -1
}

// Sys implements fs.FileInfo.
func (j *JPkgDir) Sys() any {
	return nil
}
