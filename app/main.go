package main

import (
	"os"
	"strings"

	"github.com/j4d3blooded/jpkg"
)

func main() {
	f, err := os.Create("output.dat")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := jpkg.NewJPkgEncoder(f)
	p.Name = "Test Package"

	p.AddFile(jpkg.JPkgFileToEncode{
		Source:     strings.NewReader("this is a test"),
		UUID:       jpkg.NewUUIDV4(),
		Identifier: "test1",
		Path:       "./test.txt",
		Metadata:   nil,
	})

	p.AddFile(jpkg.JPkgFileToEncode{
		Source:     strings.NewReader("this is _NOT_ a test"),
		UUID:       jpkg.NewUUIDV4(),
		Identifier: "test2",
		Path:       "./important.txt",
		Metadata: FileInfo{
			XYZ: "Additional Information",
		},
	})

	if err := p.Encode(); err != nil {
		panic(err)
	}
}

type FileInfo struct {
	XYZ string
}
