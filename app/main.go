package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/j4d3blooded/jpkg"
	jpkg_impl "github.com/j4d3blooded/jpkg/impl"
)

var (
	MODE      string
	PACKAGE   string
	DIRECTORY string
	VALID     bool
)

func init() {
	flag.StringVar(&MODE, "mode", "?", "Package mode (Pack, Unpack, Query)")
	flag.StringVar(&PACKAGE, "package", ".", "Package to unpack / output too")
	flag.StringVar(&DIRECTORY, "directory", "package.jpkg", "Directory to pack / output too")
	flag.Parse()
	MODE = strings.ToLower(MODE)
}

func main() {
	switch MODE {
	case "pack":
		pack()
	case "unpack":
		unpack()
	case "query":
		query()
	default:
		flag.PrintDefaults()
	}
}

func pack() {
	f, err := os.Create(PACKAGE)
	if err != nil {
		panic(fmt.Errorf("error creating package file: %w", err))
	}

	p := jpkg.NewJPkgEncoder(f)
	p.Name = "Archive"
	p.Compression = &jpkg_impl.LZWCompressionHandler{}

	fs.WalkDir(
		os.DirFS(DIRECTORY),
		".",
		func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			fullPath := filepath.Join(DIRECTORY, path)

			fmt.Printf("Adding file %v to package\n", fullPath)
			f2, err := os.Open(fullPath)

			if err != nil {
				panic(fmt.Errorf("could not open file %v for reading: %w", fullPath, err))
			}

			p.AddFile(
				jpkg.JPkgFileToEncode{
					Source:     f2,
					UUID:       jpkg.NewUUIDV4(),
					Identifier: "no-ident",
					Path:       path,
					Metadata:   nil,
				},
			)

			return nil
		},
	)

	if err := p.Encode(); err != nil {
		panic(fmt.Errorf("error encoding package: %w", err))
	}
}

func unpack() {
	f, err := os.Open(PACKAGE)
	if err != nil {
		panic(fmt.Errorf("error opening package: %w", err))
	}
	defer f.Close()

	pkg, err := jpkg.ReadJPkg(f, nil)
	if err != nil {
		panic(fmt.Errorf("error reading jpkg: %w", err))
	}

	fs.WalkDir(
		pkg, ".",
		func(path string, d fs.DirEntry, err error) error {

			fullPath := filepath.Join(DIRECTORY, path)

			if d.IsDir() {
				if err := os.Mkdir(fullPath, os.ModeDir); !(err == nil || errors.Is(err, os.ErrExist)) {
					panic(fmt.Errorf("error creating directory %v: %w", fullPath, err))
				}
				return nil
			}

			f1, err := pkg.GetByPath(path)
			if err != nil {
				panic(fmt.Errorf("error opening package file %v: %w", path, err))
			}

			f2, err := os.Create(fullPath)

			if err != nil {
				panic(fmt.Errorf("error creating output file %v: %w", fullPath, err))
			}
			defer f2.Close()

			if _, err := f2.ReadFrom(f1); err != nil {
				panic(fmt.Errorf("error writing to output file"))
			}

			fmt.Printf("Exporting file %v\n", fullPath)

			return nil
		},
	)
}

func query() {
	f, err := os.Open(PACKAGE)
	if err != nil {
		panic(fmt.Errorf("error opening package: %w", err))
	}
	defer f.Close()

	pkg, err := jpkg.ReadJPkg(f, nil)
	if err != nil {
		panic(fmt.Errorf("error reading jpkg: %w", err))
	}

	cFlag, eFlag := pkg.GetFlagsAndInfo()

	fmt.Printf(
		"Package Name: %v. File Count: %v. Packaged at: %v. Compressed %v, Encrypted %v\n",
		pkg.GetName(), pkg.GetPackagedTime(), pkg.GetFileCount(), cFlag != jpkg_impl.COMPRESSION_NONE, eFlag != jpkg_impl.ENCRYPTION_NONE,
	)
}
