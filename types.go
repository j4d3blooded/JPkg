package jpkg

import (
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

type JPkgFileRecord struct {
	JPkgFileRecordWithoutData
	CompressedData []byte
}

type JPkgFooter struct{}

type JPkg struct {
	Header      JPkgHeader
	Manifest    JPkgManifest
	FileRecords []JPkgFileRecord
	Footer      JPkgFooter
}
