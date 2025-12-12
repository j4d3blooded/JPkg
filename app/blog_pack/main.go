package main

import (
	"fmt"
	"io/fs"
	"os"

	jpkg "github.com/j4d3blooded/JPkg"
)

func main() {

	args := os.Args[1:]

	if len(args) < 2 || args[0] == "-h" {
		fmt.Println("HELP: blogpack [package_file] [glob_patterns]...")
		return
	}

	packageFile := args[0]
	patterns := args[1:]

	f, err := os.Create(packageFile)
	if err != nil {
		fmt.Printf("Error creating package file: %v\n", err)
		return
	}
	defer f.Close()

	pkgBuilder := jpkg.NewJPkgEncoder(f)
	fsDir := os.DirFS(".")

	for i, pattern := range patterns {
		matches, err := fs.Glob(fsDir, pattern)
		if err != nil {
			fmt.Printf("Error matching pattern %v (%v): %v\n", i, pattern, err)
			return
		}

		for _, match := range matches {

			fmt.Printf("Adding file %v\n", match)

			matchedFile, err := os.Open(match)
			if err != nil {
				fmt.Printf("Error opening file %v: %v\n", match, err)
				return
			}

			info, err := matchedFile.Stat()
			if err != nil {
				fmt.Printf("Error getting file info %v: %v\n", match, err)
				return
			}

			file := jpkg.JPkgFileToEncode{
				Source:     matchedFile,
				UUID:       jpkg.NewUUIDV4(),
				Identifier: "",
				Path:       match,
				Metadata: BlogPostMetadata{
					LastModified: info.ModTime().Unix(),
				},
			}

			if err := pkgBuilder.AddFile(file); err != nil {
				fmt.Printf("Error adding file %v: %v\n", match, err)
				return
			}
		}
	}

	if err := pkgBuilder.Encode(); err != nil {
		fmt.Printf("Error encoding package: %v\n", err)
		return
	}
}

type BlogPostMetadata struct {
	LastModified int64
}
