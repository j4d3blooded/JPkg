package jpkg

import (
	"fmt"
	"io"
	"strings"
	"time"

	jpkg_bin "github.com/j4d3blooded/jpkg/bin"
	jpkg_fs "github.com/j4d3blooded/jpkg/fs"
	jpkg_impl "github.com/j4d3blooded/jpkg/impl"
)

func ReadJPkg(r io.ReadSeeker, encryptionKey []byte) (*JPkg, error) {

	header, err := parseHeader(r)
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	manifest, err := parseManifest(r)
	if err != nil {
		return nil, fmt.Errorf("error reading manifest: %w", err)
	}

	files, err := parseFiles(r, manifest.FileCount)
	if err != nil {
		return nil, fmt.Errorf("error reading file records: %w", err)
	}

	pkg := &JPkg{
		reader:         r,
		cHandler:       jpkg_impl.GetCompressionHandler(header.CompressionFlag),
		eHandler:       jpkg_impl.GetEncryptionHandler(header.EncryptionFlag, encryptionKey),
		signatureValid: false,
		integrityValid: false,
		packagedAt:     time.Unix(manifest.PackagedAt, 0),
		name:           manifest.PackageName,
		metadata:       []byte(manifest.PackageMetadataJSON),
	}

	fileOpeners, directoryOpeners, err := buildFS(files)
	if err != nil {
		return nil, fmt.Errorf("error building file system: %w", err)
	}

	pkg.pathsToFiles = fileOpeners
	pkg.pathsToDirectories = directoryOpeners

	return pkg, nil
}

func parseHeader(r io.ReadSeeker) (*JPkgHeader, error) {
	header, err := jpkg_bin.BinaryRead[JPkgHeader](r)

	if err != nil {
		return nil, fmt.Errorf("error reading jpkg header: %w", err)
	}

	if header.MagicNumber != MAGIC_NUMBER {
		return nil, fmt.Errorf("incorrect magic number")
	}

	if header.Version != 0 {
		return nil, fmt.Errorf("unsupported version: %v", header.Version)
	}

	return header, nil
}

func parseManifest(r io.ReadSeeker) (*JPkgManifest, error) {
	header, err := jpkg_bin.BinaryRead[JPkgManifest](r)

	if err != nil {
		return nil, fmt.Errorf("error reading jpkg manifest: %w", err)
	}

	return header, nil
}

func parseFiles(r io.ReadSeeker, fileCount uint64) ([]JPkgFileRecordWithOffset, error) {
	files := make([]JPkgFileRecordWithOffset, fileCount)

	for i := range fileCount {
		record, err := jpkg_bin.BinaryRead[JPkgFileRecordWithoutData](r)
		if err != nil {
			return nil, fmt.Errorf("error reading file record %v: %w", i, err)
		}

		offset, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("error seeking in file: %w", err)
		}

		_, err = r.Seek(int64(record.CompressedDataSize), io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf("error seeking in file: %w", err)
		}

		files[i] = JPkgFileRecordWithOffset{
			JPkgFileRecordWithoutData: *record,
			Offset:                    uint64(offset),
		}
	}

	return files, nil
}

func buildFS(files []JPkgFileRecordWithOffset) (map[string]jpkgFileOpenerInfo, map[string]jpkgDirOpenerInfo, error) {
	paths := map[string]JPkgFileRecordWithOffset{}

	for _, file := range files {
		if _, exists := paths[file.FilePath]; exists {
			return nil, nil, fmt.Errorf("filepath %v is in use more then once", file.FilePath)
		}
		paths[file.FilePath] = file
	}

	treeRoot := jpkg_fs.JPkgFS{
		Root: &jpkg_fs.JPkgFSDirectory{
			Name: "\\",
		},
	}

	pathToNode := map[string]jpkg_fs.JPkgFSNode{
		"\\": treeRoot.Root,
	}

	for path := range paths {
		segments := strings.Split(path, "\\")
		if len(segments) == 0 || len(segments) == 1 {
			return nil, nil, fmt.Errorf("malformed path: %v", path)
		}
		dirSegs := segments[1 : len(segments)-1]
		lastSeg := segments[len(segments)-1]

		var current jpkg_fs.JPkgFSNode = treeRoot.Root

	outer:
		for _, seg := range dirSegs {
			dir, isDir := current.(*jpkg_fs.JPkgFSDirectory)
			if !isDir {
				fullPath := jpkg_fs.GetFullPath(current)
				return nil, nil, fmt.Errorf("path %v is a file but is being used as a directory for %v", fullPath, path)
			}

			for _, child := range dir.Children {
				if child.GetName() == seg {
					current = child
					continue outer
				}
			}

			next := &jpkg_fs.JPkgFSDirectory{
				Parent:   dir,
				Name:     seg,
				Children: nil,
			}
			pathToNode[jpkg_fs.GetFullPath(next)] = next
			dir.Children = append(dir.Children, next)
			current = next
		}

		dir, isDir := current.(*jpkg_fs.JPkgFSDirectory)
		if !isDir {
			fullPath := jpkg_fs.GetFullPath(current)
			return nil, nil, fmt.Errorf("path %v is used as a file but is a directory for %v", fullPath, path)
		}

		for _, child := range dir.Children {
			if child.GetName() == lastSeg {
				return nil, nil, fmt.Errorf("filename of path is already in use as either directory or file: %v", path)
			}
		}

		file := &jpkg_fs.JPkgFSFile{
			Parent: dir,
			Name:   lastSeg,
		}

		dir.Children = append(dir.Children, file)
		pathToNode[path] = file
	}

	fils := map[string]jpkgFileOpenerInfo{}
	dirs := map[string]jpkgDirOpenerInfo{}

	for path, node := range pathToNode {
		convertNodeToOpenerInfo(node, path, dirs, fils, paths)
	}

	return fils, dirs, nil
}

func convertNodeToOpenerInfo(
	node jpkg_fs.JPkgFSNode, path string,
	directories map[string]jpkgDirOpenerInfo, files map[string]jpkgFileOpenerInfo,
	paths map[string]JPkgFileRecordWithOffset,
) {
	switch f := node.(type) {
	case *jpkg_fs.JPkgFSDirectory:

		childPaths := make([]string, len(f.Children))
		for i, v := range f.Children {
			childPaths[i] = jpkg_fs.GetFullPath(v)
		}

		directories[path] = jpkgDirOpenerInfo{
			name:       f.Name,
			path:       jpkg_fs.GetFullPath(f),
			ChildPaths: childPaths,
		}

	case *jpkg_fs.JPkgFSFile:
		files[path] = jpkgFileOpenerInfo{
			name:             f.Name,
			compressedSize:   paths[path].CompressedDataSize,
			uncompressedSize: paths[path].UncompressedDataSize,
			path:             paths[path].FilePath,
			identifier:       paths[path].FileIdentifier,
			uuid:             paths[path].UUID,
			offset:           int64(paths[path].Offset),
			metadata:         []byte(paths[path].FileMetadataJSON),
		}

	default:
		panic("unsupported fs node type")
	}
}
