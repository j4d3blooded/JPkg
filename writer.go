package jpkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
)

func NewJPkgEncoder(w io.Writer) *JPkgEncoder {
	wr := &JPkgEncoder{
		w:           w,
		Encryption:  &NullEncryptionHandler{},
		Compression: &NullCompressionHandler{},
		Hasher:      &NullHasherHandler{},
		Signer:      &NullCryptoHandler{},
		PackageTime: time.Now(),
		Metadata:    nil,
	}
	return wr
}

type JPkgEncoder struct {
	Name        string
	Metadata    any
	PackageTime time.Time
	Encryption  EncryptionHandler
	Compression CompressionHandler
	Hasher      HasherHandler
	Signer      CryptoHandler
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
	source           io.Reader
	uuid             UUID
	identifier       string
	identifierLength uint64
	metadataJson     string
	metadataLength   uint64
	path             string
	pathLength       uint64
}

func (j *JPkgEncoder) AddFile(file JPkgFileToEncode) error {

	if file.Source == nil {
		return errors.New("file has no source")
	}

	json, jsonSize, err := serializeMetadataToJSON(file.Metadata)
	if err != nil {
		return fmt.Errorf("error serializing json metadata: %w", err)
	}

	nf := jpkgFileRecord{
		source:           file.Source,
		uuid:             file.UUID,
		identifier:       file.Identifier,
		identifierLength: uint64(len(file.Identifier)),
		metadataJson:     json,
		metadataLength:   jsonSize,
		path:             file.Path,
		pathLength:       uint64(len(file.Path)),
	}

	j.files = append(j.files, nf)

	return nil
}

func (j *JPkgEncoder) Encode() error {

	o := uint64(0)
	o = j.writeHeader(o)
	o, err := j.writeManifest(o)
	if err != nil {
		return fmt.Errorf("error writing package manifest: %w", err)
	}
	o, err = j.writeBody(o)
	if err != nil {
		return fmt.Errorf("error writing package body: %w", err)
	}

	return nil
}

func (j *JPkgEncoder) writeHeader(o uint64) uint64 {
	j.ws("jpkg") // Magic Number
	o += 4

	j.wb(uint64(1)) // Version
	o += 8

	j.wb(j.Compression.Flag()) // Compression Flag
	j.wb(j.Encryption.Flag())  // Encryption Flag
	j.wb(j.Hasher.Flag())      // Hash Flag
	j.wb(j.Signer.Flag())      // Crypto Flag
	o += 4                     // flags

	j.pad(16) // Padding
	o += 16

	return o
}

func (j *JPkgEncoder) writeManifest(o uint64) (uint64, error) {
	j.wb(j.PackageTime.Unix()) // Timestamp
	o += 8

	j.wb(uint64(len(j.files))) // File Count
	o += 8

	j.wb(uint64(len(j.Name))) // Package Name Length
	o += 8

	json, jsonLen, err := serializeMetadataToJSON(j.Metadata)
	if err != nil {
		return o, fmt.Errorf("error serializing metadata: %w", err)
	}

	j.wb(jsonLen) // Package Metadata Length
	o += 8

	j.ws(j.Name) // Package Name
	o += uint64(len(j.Name))

	j.ws(json) // Package Metadata
	o += jsonLen

	padding := calculatePaddingLength(o) // Pad 16
	j.pad(padding)
	o += padding

	j.pad(16) // Padding
	o += 16

	return o, nil
}

func (j *JPkgEncoder) writeBody(o uint64) (uint64, error) {

	for _, file := range j.files {
		compressed := &bytes.Buffer{}
		compressor := j.Compression.Compress(compressed)
		uncompressedSize, err := io.Copy(compressor, file.source)

		if err != nil {
			return o, fmt.Errorf("error compressing file %v: %w", file.path, err)
		}

		if err := compressor.Close(); err != nil {
			return o, fmt.Errorf("error finalizing compression for file %v: %w", file.path, err)
		}

		encrypted := &bytes.Buffer{}
		encryptor, err := j.Encryption.Encrypt(encrypted)
		if err != nil {
			return o, fmt.Errorf("error creating encryptor: %w", err)
		}

		if _, err := compressed.WriteTo(encryptor); err != nil {
			return o, fmt.Errorf("error encrypting file %v: %w", file.path, err)
		}
		if err := encryptor.Close(); err != nil {
			return o, fmt.Errorf("error finalizing encryption for file %v: %w", file.path, err)
		}

		compressedSize := uint64(encrypted.Len())

		fw := newBufferWriter()

		fw.wb(file.identifierLength) // identifier length
		fw.wb(file.pathLength)       // path length
		fw.wb(file.uuid)             // uuid
		fw.wb(file.metadataLength)   // metadata length
		fw.wb(compressedSize)        // compressed data size
		fw.wb(uncompressedSize)      // uncompressed data size
		fw.ws(file.identifier)       // identifier
		fw.ws(file.path)             // path
		fw.ws(file.metadataJson)     // metadata json
		fw.b.ReadFrom(encrypted)     // compressed data

		blockSize := fw.b.Len()
		blockSize += 8

		o += uint64(blockSize)

		padding := calculatePaddingLength(o)
		fw.pad(padding)
		o += padding
		j.wb(o)

		fw.b.WriteTo(j.w)
	}

	j.pad(16)
	o += 16

	return o, nil
}

func (j *JPkgEncoder) wb(data any) {
	binary.Write(j.w, binary.BigEndian, data)
}

func (j *JPkgEncoder) ws(data string) {
	io.WriteString(j.w, data)
}

func (j *JPkgEncoder) pad(count uint64) {
	db := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	for i := range int(count) {
		n := i % len(db)
		j.wb(db[n])
	}
}

func calculatePaddingLength(M uint64) uint64 {
	padding := M / 16
	padding = padding + 1
	padding = 16 * padding
	padding = padding - M
	padding = padding % uint64(16)
	return padding
}
