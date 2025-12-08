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

	w := jpkg.NewJPkgEncoder(f)
	w.AddFile("test.txt", strings.NewReader("this is a test"), FileInfo{"XYZ Info"})
	w.Encode()
}

type FileInfo struct {
	XYZ string
}
