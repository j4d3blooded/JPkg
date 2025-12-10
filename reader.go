package jpkg

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"

	jpkg_impl "github.com/j4d3blooded/jpkg/impl"
)

func NewJPkgDecoder[T any](r io.ReadSeeker, encKey []byte) (*JPkgDecoder[T], error) {
	rr := &wrapReader{r}

	compression, encryption, hash, signature, err := parseHeader(rr, encKey)
	if err != nil {
		return nil, err
	}
	pkgTime, fileCount, pkgName, pkgMetadata, offset, err := parseManifest[T](rr)
	if err != nil {
		return nil, err
	}

	rr.Seek(int64(offset), io.SeekStart)

	files := parseFileRecords(fileCount, rr, offset)

	return &JPkgDecoder{
		pkgName,
		compression,
		encryption,
		hash,
		signature,
		*pkgMetadata,
		fileCount,
		pkgTime,
		rr,
		files,
	}, nil
}

type JPkgDecoder struct {
	Name        string
	Compression jpkg_impl.CompressionHandler
	Encryption  jpkg_impl.EncryptionHandler
	Hasher      jpkg_impl.HasherHandler
	Crypto      jpkg_impl.CryptoHandler
	Metadata    []byte // json byte array
	FileCount   uint64
	PackageTime time.Time
	r           *wrapReader
	files       []fileReadInfo
}

type fileReadInfo struct {
	dataOffset  uint64
	dataSize    uint64
	decompSize  uint64
	identifier  string
	uuid        UUID
	path        string
	metadataStr string
}

func parseHeader(rr *wrapReader, encKey []byte) (jpkg_impl.CompressionHandler, jpkg_impl.EncryptionHandler, jpkg_impl.HasherHandler, jpkg_impl.CryptoHandler, error) {
	if s := rr.readStr(4); s != "jpkg" {
		return nil, nil, nil, nil, errors.New("invalid magic number")
	}

	if version := rr.u64(); version != 1 {
		return nil, nil, nil, nil, errors.New("unsupported version")
	}

	compression := jpkg_impl.GetCompressionHandler(jpkg_impl.CompressionFlag(rr.u8()))
	encryption := jpkg_impl.GetEncryptionHandler(jpkg_impl.EncryptionFlag(rr.u8()), encKey)
	hash := jpkg_impl.GetHashHandler(jpkg_impl.HasherFlag(rr.u8()))
	signature := jpkg_impl.GetCryptoHandler(jpkg_impl.CryptoFlag(rr.u8()))

	rr.readN(16)
	return compression, encryption, hash, signature, nil
}

func parseManifest(rr *wrapReader) (time.Time, uint64, string, []byte, uint64, error) {
	pkgTime := time.Unix(int64(rr.u64()), 0)

	fileCount := rr.u64()

	pkgNameLength := rr.u64()
	pkgMetadataLength := rr.u64()

	pkgName := rr.readStr(pkgNameLength)
	pkgMetadataJson := rr.readStr(pkgMetadataLength)

	offset := uint64(0x20 + 0x20)
	offset += pkgNameLength + pkgMetadataLength
	offset += calculatePaddingLength(pkgNameLength + pkgMetadataLength)
	offset += 0x10

	return pkgTime, fileCount, pkgName, []byte(pkgMetadataJson), offset, nil
}

func parseFileRecords(fileCount uint64, rr *wrapReader, nextOffset uint64) []fileReadInfo {
	files := []fileReadInfo{}

	for range fileCount {
		rr.Seek(int64(nextOffset), io.SeekStart)
		recordStart := nextOffset

		nextOffset = rr.u64()
		idenLength := rr.u64()
		pathLength := rr.u64()
		uuid := readT[UUID](rr)
		mdLength := rr.u64()
		cSize := rr.u64()
		dSize := rr.u64()

		ident := rr.readStr(idenLength)
		path := rr.readStr(pathLength)
		meta := rr.readStr(mdLength)

		files = append(files,
			fileReadInfo{
				dataOffset:  recordStart + 0x40 + idenLength + pathLength + mdLength,
				dataSize:    cSize,
				identifier:  ident,
				uuid:        uuid,
				path:        path,
				metadataStr: meta,
				decompSize:  dSize,
			},
		)
	}
	return files
}

// Open implements fs.FS.
func (j *JPkgDecoder) Open(name string) (fs.File, error) {
	panic("unimplemented")
}

func buildFSTree(records []jpkgFileRecord) fsTree {

	root := JPkgFSDirectory{path: "/"}

	for _, file := range records {
		current := &root
		segments := strings.Split(file.path, "\\\\")
		for _, segment := range segments[:len(segments)-1] {

		}
	}

	return root
}
