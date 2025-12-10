package main

import (
	"os"
	"strings"

	"github.com/j4d3blooded/jpkg"
)

func main() {
	main_write()
	main_read()
}

func main_write() {
	f, err := os.Create("output.dat")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := jpkg.NewJPkgEncoder(f)
	p.Name = "Test Package"

	p.Metadata = FileInfo{
		XYZ: "This is a test package metadata",
	}

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

	p.AddFile(jpkg.JPkgFileToEncode{
		Source:     strings.NewReader("this is _NOT_ a test"),
		UUID:       jpkg.NewUUIDV4(),
		Identifier: "test2",
		Path:       "./f2/important.txt",
		Metadata: FileInfo{
			XYZ: "Additional Information",
		},
	})

	p.AddFile(jpkg.JPkgFileToEncode{
		Source:     strings.NewReader("this is _NOT_ a test"),
		UUID:       jpkg.NewUUIDV4(),
		Identifier: "test2",
		Path:       "./f2/important_2.txt",
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

func main_read() {
	f, err := os.Open("output.dat")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = jpkg.ReadJPkg(f, nil)

	if err != nil {
		panic(err)
	}

}
