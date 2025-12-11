package jpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"time"

	jpkg_impl "github.com/j4d3blooded/jpkg/impl"
)

type JPkgHeader struct {
	MagicNumber     uint32
	Version         uint64
	CompressionFlag jpkg_impl.CompressionFlag
	EncryptionFlag  jpkg_impl.EncryptionFlag
	HasherFlag      jpkg_impl.HasherFlag
	SignatureFlag   jpkg_impl.CryptoFlag
}

type JPkgManifest struct {
	PackagedAt          int64
	FileCount           uint64
	PackageName         string
	PackageMetadataJSON string
}

type JPkgFileRecordWithoutData struct {
	FileIdentifier       string
	FilePath             string
	UUID                 UUID
	FileMetadataJSON     string
	CompressedDataSize   uint64
	UncompressedDataSize uint64
}

type _JPkgFileRecord struct {
	JPkgFileRecordWithoutData
	CompressedData []byte
}

type JPkgFileRecordWithOffset struct {
	JPkgFileRecordWithoutData
	Offset uint64
}

type JPkg struct {
	reader             io.ReadSeeker
	pathsToFiles       map[string]jpkgFileOpenerInfo
	pathsToDirectories map[string]jpkgDirOpenerInfo
	cHandler           jpkg_impl.CompressionHandler
	eHandler           jpkg_impl.EncryptionHandler
	signatureValid     bool
	integrityValid     bool
	packagedAt         time.Time
	name               string
	metadata           []byte
}

func (j *JPkg) Open(name string) (fs.File, error) {

	name = normalizeFilePath(name)

	if fileInfo, isFile := j.pathsToFiles[name]; isFile {
		_, err := j.reader.Seek(fileInfo.offset, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("error seeking with package reader: %w", err)
		}

		decrypted := bytes.Buffer{}
		decryptor, err := j.eHandler.Decrypt(&decrypted)
		if err != nil {
			return nil, fmt.Errorf("error creating aes decryptor: %w", err)
		}
		if _, err := io.CopyN(decryptor, j.reader, int64(fileInfo.compressedSize)); err != nil {
			return nil, fmt.Errorf("error decrypting file data: %w", err)
		}

		decompressed, err := j.cHandler.Decompress(decrypted.Bytes())
		if err != nil {
			return nil, fmt.Errorf("error decompressing file data: %w", err)
		}

		buffer := bytes.NewReader(decompressed)

		return &JPkgFile{
			pkg:        j,
			name:       fileInfo.name,
			size:       int64(fileInfo.uncompressedSize),
			closed:     false,
			buffer:     *buffer,
			path:       fileInfo.path,
			identifier: fileInfo.identifier,
			uuid:       fileInfo.uuid,
		}, nil
	}

	if dirInfo, isDir := j.pathsToDirectories[name]; isDir {
		return &JPkgDir{
			pkg:    j,
			name:   dirInfo.name,
			path:   dirInfo.path,
			dirIdx: 0,
		}, nil
	}

	return nil, fs.ErrNotExist
}

func (j *JPkg) ReadDir(name string) ([]fs.DirEntry, error) {

	name = normalizeFilePath(name)

	dirInfo, isDir := j.pathsToDirectories[name]

	if !isDir {
		return nil, fs.ErrNotExist
	}

	entries := make([]fs.DirEntry, len(dirInfo.ChildPaths))

	for i, child := range dirInfo.ChildPaths {
		if dirInfo, isDir := j.pathsToDirectories[child]; isDir {
			entries[i] = &JPkgDirInfo{
				pkg:   j,
				path:  child,
				name:  dirInfo.name,
				size:  0,
				isDir: true,
			}
		} else if fileInfo, isFile := j.pathsToFiles[child]; isFile {
			entries[i] = &JPkgDirInfo{
				pkg:   j,
				path:  child,
				name:  fileInfo.name,
				size:  int64(fileInfo.uncompressedSize),
				isDir: false,
			}
		} else {
			panic(fmt.Errorf("child is not real? %v", child))
		}
	}

	return entries, nil
}

func (j *JPkg) GetName() string {
	return j.name
}

func GetPackageMetadata[T any](pkg *JPkg) (*T, error) {
	v := new(T)
	err := json.Unmarshal(pkg.metadata, v)
	if err != nil {
		return nil, fmt.Errorf("error parsing package metadata: %w", err)
	}
	return v, nil
}

func (j *JPkg) GetByUUID(uuid UUID) (*JPkgFile, error) {
	for _, file := range j.pathsToFiles {
		if file.uuid == uuid {
			f, err := j.Open(file.path)
			if err != nil {
				return nil, fmt.Errorf("error opening found: %w", err)
			}
			return f.(*JPkgFile), nil
		}
	}
	return nil, errors.New("uuid not found")
}

func (j *JPkg) GetByIdentifier(expr *regexp.Regexp) ([]*JPkgFile, error) {
	files := []*JPkgFile{}

	for _, file := range j.pathsToFiles {
		if expr.MatchString(file.identifier) {
			file, err := j.Open(file.path)
			if err != nil {
				return nil, fmt.Errorf("error opening matched file: %w", err)
			}
			files = append(files, file.(*JPkgFile))
		}
	}

	return files, nil
}

func GetByMetadataQuery[T any](j *JPkg, query func(md T) bool) ([]*JPkgFile, error) {
	files := []*JPkgFile{}

	for _, file := range j.pathsToFiles {
		metadata := new(T)
		if err := json.Unmarshal(file.metadata, metadata); err != nil {
			return nil, fmt.Errorf("error parsing file metadata: %w", err)
		}
		if query(*metadata) {
			file, err := j.Open(file.path)
			if err != nil {
				return nil, fmt.Errorf("error opening matched file: %w", err)
			}
			files = append(files, file.(*JPkgFile))
		}
	}

	return files, nil
}

func (j *JPkg) GetByPath(path string) (*JPkgFile, error) {
	f, err := j.Open(path)
	if err != nil {
		return nil, err
	}
	return f.(*JPkgFile), nil
}
