package jpkg

import (
	"bytes"
	"io/fs"
	"time"
)

type jpkgFileOpenerInfo struct {
	name             string
	path             string
	identifier       string
	uuid             UUID
	compressedSize   uint64
	uncompressedSize uint64
	offset           int64
}

type JPkgFile struct {
	pkg        *JPkg
	name       string
	path       string
	identifier string
	uuid       UUID
	size       int64
	buffer     bytes.Reader
	closed     bool
}

func (j *JPkgFile) IsDir() bool {
	return false
}

func (j *JPkgFile) ModTime() time.Time {
	return j.pkg.packagedAt
}

func (j *JPkgFile) Mode() fs.FileMode {
	return 4
}

func (j *JPkgFile) Name() string {
	return j.name
}

func (j *JPkgFile) Size() int64 {
	return j.size
}

func (j *JPkgFile) Read(b []byte) (int, error) {
	if j.closed {
		return -1, fs.ErrClosed
	}

	return j.buffer.Read(b)
}

func (j *JPkgFile) Stat() (fs.FileInfo, error) {
	return j, nil
}

func (j *JPkgFile) Close() error {
	j.closed = true
	return nil
}

func (j *JPkgFile) Sys() any {
	return nil
}
