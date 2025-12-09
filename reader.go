package jpkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	jpkg_impl "github.com/j4d3blooded/jpkg/impl"
)

func NewJPkgDecoder[T any](r io.ReadSeeker, encKey []byte) (*JPkgDecoder[T], error) {
	rr := &wrapReader{r}

	if s := rr.readStr(4); s != "jpkg" {
		return nil, errors.New("invalid magic number")
	}

	if version := rr.u64(); version != 1 {
		return nil, errors.New("unsupported version")
	}

	compression := jpkg_impl.GetCompressionHandler(jpkg_impl.CompressionFlag(rr.u8()))
	encryption := jpkg_impl.GetEncryptionHandler(jpkg_impl.EncryptionFlag(rr.u8()), encKey)
	hash := jpkg_impl.GetHashHandler(jpkg_impl.HasherFlag(rr.u8()))
	signature := jpkg_impl.GetCryptoHandler(jpkg_impl.CryptoFlag(rr.u8()))

	rr.readN(16)
	pkgTime := time.Unix(int64(rr.u64()), 0)

	fileCount := rr.u64()

	pkgNameLength := rr.u64()
	pkgMetadataLength := rr.u64()

	pkgName := rr.readStr(pkgNameLength)
	pkgMetadataJson := rr.readStr(pkgMetadataLength)

	pkgMetadata := new(T)
	if err := json.Unmarshal([]byte(pkgMetadataJson), pkgMetadata); err != nil {
		return nil, fmt.Errorf("error parsing package metadata: %w", err)
	}

	offset := 0x20 + 0x20 + pkgNameLength + pkgMetadataLength + calculatePaddingLength(pkgNameLength+pkgMetadataLength) + 0x10

	rr.Seek(int64(offset), io.SeekStart)

	files := readFiles(fileCount, rr, offset)

	return &JPkgDecoder[T]{
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

func readFiles(fileCount uint64, rr *wrapReader, nextOffset uint64) []fileReadInfo {
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

type JPkgDecoder[T any] struct {
	Name        string
	Compression jpkg_impl.CompressionHandler
	Encryption  jpkg_impl.EncryptionHandler
	Hasher      jpkg_impl.HasherHandler
	Crypto      jpkg_impl.CryptoHandler
	Metadata    T
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
