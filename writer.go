package jpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	jpkg_bin "github.com/j4d3blooded/JPkg/bin"
	jpkg_impl "github.com/j4d3blooded/JPkg/impl"
)

func NewJPkgEncoder(w io.Writer) *JPkgEncoder {
	wr := &JPkgEncoder{
		w:           w,
		Encryption:  &jpkg_impl.NullEncryptionHandler{},
		Compression: &jpkg_impl.NullCompressionHandler{},
		Hasher:      &jpkg_impl.NullHasherHandler{},
		Signer:      &jpkg_impl.NullCryptoHandler{},
		PackageTime: time.Now(),
		Metadata:    nil,
	}
	return wr
}

type JPkgEncoder struct {
	Name        string
	Metadata    any
	PackageTime time.Time
	Encryption  jpkg_impl.EncryptionHandler
	Compression jpkg_impl.CompressionHandler
	Hasher      jpkg_impl.HasherHandler
	Signer      jpkg_impl.CryptoHandler
	w           io.Writer
	files       []jpkgFileRecord
	offset      uint64
}

type JPkgFileToEncode struct {
	Source     io.Reader
	UUID       UUID
	Identifier string
	Path       string
	Metadata   any
}

type jpkgFileRecord struct {
	source       io.Reader
	identifier   string
	path         string
	uuid         UUID
	metadataJson string
}

func (j *JPkgEncoder) AddFile(file JPkgFileToEncode) error {

	if file.Source == nil {
		return errors.New("file has no source")
	}

	json, err := serializeMetadataToJSON(file.Metadata)
	if err != nil {
		return fmt.Errorf("error serializing json metadata: %w", err)
	}

	file.Path = strings.TrimPrefix(file.Path, ".")

	if !strings.HasPrefix(file.Path, "/") {
		file.Path = "/" + file.Path
	}

	file.Path = filepath.Clean(file.Path)

	for _, existingFile := range j.files {
		if existingFile.path == file.Path {
			return errors.New("path is already in use")
		}
	}

	nf := jpkgFileRecord{
		source:       file.Source,
		uuid:         file.UUID,
		identifier:   file.Identifier,
		metadataJson: json,
		path:         file.Path,
	}

	j.files = append(j.files, nf)

	return nil
}

func (j *JPkgEncoder) Encode() error {
	if err := j.writeHeader(); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	if err := j.writeManifest(); err != nil {
		return fmt.Errorf("error writing manifest: %w", err)
	}

	if err := j.writeFileRecords(); err != nil {
		return fmt.Errorf("error writing file records: %w", err)
	}

	return nil
}

func (j *JPkgEncoder) writeHeader() error {
	header := JPkgHeader{
		MagicNumber:     MAGIC_NUMBER,
		Version:         0,
		CompressionFlag: j.Compression.Flag(),
		EncryptionFlag:  j.Encryption.Flag(),
		HasherFlag:      j.Hasher.Flag(),
		SignatureFlag:   j.Signer.Flag(),
	}

	return jpkg_bin.BinaryWrite(j.w, header)
}

func (j *JPkgEncoder) writeManifest() error {

	metadataJson, err := json.Marshal(j.Metadata)
	if err != nil {
		return fmt.Errorf("error converting package metadata to json: %w", err)
	}

	manifest := JPkgManifest{
		PackagedAt:          j.PackageTime.Unix(),
		FileCount:           uint64(len(j.files)),
		PackageName:         j.Name,
		PackageMetadataJSON: string(metadataJson),
	}

	if err := jpkg_bin.BinaryWrite(j.w, manifest); err != nil {
		return fmt.Errorf("error writing package manifest: %w", err)
	}

	return nil
}

func (j *JPkgEncoder) writeFileRecords() error {
	for _, file := range j.files {

		uncompressedBytes, err := io.ReadAll(file.source)
		if err != nil {
			return fmt.Errorf("error reading file (%v/%v/%v): %w", file.path, file.identifier, file.uuid, err)
		}

		uncompressedSize := len(uncompressedBytes)
		compressed, err := j.Compression.Compress(uncompressedBytes)
		if err != nil {
			return fmt.Errorf("error compressing file (%v/%v/%v): %w", file.path, file.identifier, file.uuid, err)
		}

		encrypted := bytes.Buffer{}
		encryptor, err := j.Encryption.Encrypt(&encrypted)
		if err != nil {
			return fmt.Errorf("error creating encryptor: %w", err)
		}
		if _, err := encryptor.Write(compressed); err != nil {
			return fmt.Errorf("error encrypting file (%v/%v/%v): %w", file.path, file.identifier, file.uuid, err)
		}

		record := _JPkgFileRecord{
			JPkgFileRecordWithoutData: JPkgFileRecordWithoutData{
				FileIdentifier:       file.identifier,
				FilePath:             file.path,
				UUID:                 file.uuid,
				FileMetadataJSON:     file.metadataJson,
				CompressedDataSize:   uint64(encrypted.Len()),
				UncompressedDataSize: uint64(uncompressedSize),
			},
			CompressedData: encrypted.Bytes(),
		}

		if err := jpkg_bin.BinaryWrite(j.w, record); err != nil {
			return fmt.Errorf("error writing file (%v/%v/%v): %w", file.path, file.identifier, file.uuid, err)
		}
	}

	return nil
}
