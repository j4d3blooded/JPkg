package jpkg

import (
	"io/fs"
	"time"
)

type JPkgDirInfo struct {
	pkg   *JPkg
	path  string
	name  string
	size  int64
	isDir bool
}

// Info implements fs.DirEntry.
func (j *JPkgDirInfo) Info() (fs.FileInfo, error) {
	return j, nil
}

// IsDir implements fs.DirEntry.
func (j *JPkgDirInfo) IsDir() bool {
	_, isDir := j.pkg.pathsToDirectories[j.path]
	return isDir
}

// Name implements fs.DirEntry.
func (j *JPkgDirInfo) Name() string {
	return j.name
}

// Type implements fs.DirEntry.
func (j *JPkgDirInfo) Type() fs.FileMode {
	if j.isDir {
		return fs.ModeDir
	}
	return 4
}

// ModTime implements fs.FileInfo.
func (j *JPkgDirInfo) ModTime() time.Time {
	return j.pkg.packagedAt
}

// Mode implements fs.FileInfo.
func (j *JPkgDirInfo) Mode() fs.FileMode {
	return j.Type()
}

// Size implements fs.FileInfo.
func (j *JPkgDirInfo) Size() int64 {
	return j.size
}

// Sys implements fs.FileInfo.
func (j *JPkgDirInfo) Sys() any {
	return nil
}
