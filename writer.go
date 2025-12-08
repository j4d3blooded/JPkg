package jpkg

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"time"
)

func NewJPkgEncoder(w io.Writer) *JPkgEncoder {
	wr := &JPkgEncoder{
		w:           w,
		Encryption:  &NullEncryptionHandler{},
		Compression: &NullCompressionHandler{},
		PackageTime: time.Now(),
	}
	return wr
}

type JPkgEncoder struct {
	PackageTime time.Time
	Encryption  EncryptionHandler
	Compression CompressionHandler
	Metadata    any
	w           io.Writer
	files       []jPkgFileToEncode
}

func (j *JPkgEncoder) SetEncryption(e EncryptionHandler) {
	j.Encryption = e
}

func (j *JPkgEncoder) SetCompression(c CompressionHandler) {
	j.Compression = c
}

type jPkgFileToEncode struct {
	path           string
	pathLength     uint64
	metadataJson   string
	metadataLength uint64
	data           io.Reader
	dataLength     uint64
}

func (j *JPkgEncoder) AddFile(path string, source io.Reader, metadata any) error {

	json, jsonLen, err := serializeMetadataToJSON(j.Metadata)
	if err != nil {
		return err
	}

	file := jPkgFileToEncode{
		path:           path,
		pathLength:     uint64(len(path)),
		metadataJson:   json,
		metadataLength: jsonLen,
	}

	if size, isSizable := isSizeableReader(source); isSizable {
		file.dataLength = size
		file.data = source
		j.files = append(j.files, file)
		return nil
	}

	if f, isFile := source.(fs.File); isFile {
		info, err := f.Stat()
		if err != nil {
			return fmt.Errorf("error getting file stats: %w", err)
		}
		file.data = source
		file.dataLength = uint64(info.Size())
		j.files = append(j.files, file)
		return nil
	}

	b, err := io.ReadAll(source)
	if err != nil {
		return fmt.Errorf("error reading source: %w", err)
	}

	file.data = bytes.NewReader(b)
	file.dataLength = uint64(len(b))
	j.files = append(j.files, file)

	return nil
}

func (j *JPkgEncoder) Encode(hash bool, sign *rsa.PrivateKey) error {

	// Header
	j.writeHeader()

	offset, err := j.writeManifest()
	if err != nil {
		return fmt.Errorf("error writing package manifest: %w", err)
	}

	j.writeFileBlock(offset)

	if err := j.writeBody(); err != nil {
		return fmt.Errorf("error writing package body: %w", err)
	}

	return nil
}

func (j *JPkgEncoder) writeHeader() {
	j.ws("jpkg")               // Magic Number
	j.wb(uint64(1))            // Version
	j.wb(j.Compression.Flag()) // Compression Flag
	j.wb(j.Encryption.Flag())  // Encryption Flag
	j.wb(uint16(65535))        // 2 bytes FF
	j.pad(16)
}

func (j *JPkgEncoder) writeManifest() (uint64, error) {

	json, jsonLen, err := serializeMetadataToJSON(j.Metadata)
	if err != nil {
		return 0, err
	}

	j.wb(j.PackageTime.Unix()) // packaged at time
	j.wb(uint64(len(j.files))) // file count
	j.wb(jsonLen)              // package metadata length
	j.ws(json)                 // package metadata

	jsonPaddingLength := calculatePaddingLength(jsonLen + 8)

	j.pad(jsonPaddingLength)
	j.pad(16)

	offset := uint64(0)

	offset += 0x20                        // header size
	offset += 0x18                        // time + file count + metadata length
	offset += jsonLen                     // metadata
	offset += jsonPaddingLength           // metadata padding
	offset += 0x10                        // additional padding
	offset += 0x30 * uint64(len(j.files)) // length of file block

	return offset, nil
}

func (j *JPkgEncoder) writeFileBlock(offset uint64) {
	for _, file := range j.files {
		j.wb(offset)                  // File Name Offset
		j.wb(file.pathLength)         // File Name Length
		offset += file.pathLength     // Increment Offset
		j.wb(offset)                  // Metadata Offset
		j.wb(file.metadataLength)     // Metadata Length
		offset += file.metadataLength // Increment Offset
		j.wb(offset)                  // Data Offset
		j.wb(file.dataLength)         // Data Length
		offset += file.dataLength     // Increment Offset
	}
}

func (j *JPkgEncoder) writeBody() error {
	for _, file := range j.files {
		j.ws(file.path)
		j.ws(file.metadataJson)
		io.Copy(j.w, file.data)
	}
	return nil
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
