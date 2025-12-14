package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

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

			blogpostMetadata := BlogPostMetadata{
				LastModified: info.ModTime(),
			}

			metadataFile := strings.TrimSuffix(match, filepath.Ext(match)) + ".json"

			if data, err := os.ReadFile(metadataFile); err == nil {
				if err := json.Unmarshal(data, &blogpostMetadata); err != nil {
					fmt.Printf("Error unmarshaling metadata file for %v: %v\n", match, err)
				}
				fmt.Printf("Using metadata file for %v\n", match)
			}

			file := jpkg.JPkgFileToEncode{
				Source:     matchedFile,
				UUID:       jpkg.NewUUIDV4(),
				Identifier: "",
				Path:       match,
				Metadata:   blogpostMetadata,
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
	LastModified time.Time
}
