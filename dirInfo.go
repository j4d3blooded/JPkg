package jpkg

import (
	"io/fs"
)

type JPkgDirInfo struct{}

func (j *JPkgDirInfo) Close() error {

}

func (j *JPkgDirInfo) Read([]byte) (int, error) {

}

func (j *JPkgDirInfo) Stat() (fs.FileInfo, error) {

}
