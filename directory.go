package jpkg

import (
	"io/fs"
)

type jpkgDirOpenerInfo struct {
	Children uint64
}

type JPkgDir struct{}

func (j *JPkgDir) Close() error {

}

func (j *JPkgDir) Read([]byte) (int, error) {

}

func (j *JPkgDir) Stat() (fs.FileInfo, error) {

}
