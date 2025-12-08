package jpkg

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewJPkgEncoder(t *testing.T) {
	encoder := NewJPkgEncoder(nil)
	require.WithinDuration(t, encoder.PackageTime, time.Now(), time.Second)
	require.Equal(t, &NullCompressionHandler{}, encoder.Compression)
	require.Equal(t, &NullEncryptionHandler{}, encoder.Encryption)
}

func TestJPkgEncoder_AddFile(t *testing.T) {
	encoder := NewJPkgEncoder(nil)

	path := "test.txt"
	text := "xyz"
	err := encoder.AddFile(path, strings.NewReader(text), nil)

	require.Nil(t, err)
	require.Len(t, encoder.files, 1)
	require.Equal(t, uint64(3), encoder.files[0].dataLength)

	require.Equal(t, path, encoder.files[0].path)
	require.Equal(t, uint64(len(path)), encoder.files[0].pathLength)

	require.Equal(t, "{}", encoder.files[0].metadataJson)
	require.Equal(t, uint64(2), encoder.files[0].metadataLength)
}

func TestJPkgEncoder_writeHeader(t *testing.T) {
	b := &bytes.Buffer{}

	encoder := NewJPkgEncoder(b)

	encoder.Encryption = &testEncryptionFlag{}
	encoder.Compression = &testCompression{}

	encoder.writeHeader()

	require.Equal(t, "jpkg", readStr(b, 4))                            // Magic Number
	require.Equal(t, uint64(1), readU64(b))                            // Version
	require.Equal(t, TEST_COMPRESSION_FLAG, CompressionFlag(readB(b))) // Test Compression
	require.Equal(t, TEST_ENCRYPTION_FLAG, EncryptionFlag(readB(b)))   // Test Encryption
	require.Equal(t, byte(0xFF), readB(b))
	require.Equal(t, byte(0xFF), readB(b))
}

func readStr(r io.Reader, n int) string {
	b := make([]byte, n)
	r.Read(b)
	return string(b)
}

// HELPER

func readB(r io.Reader) byte {
	var b byte
	binary.Read(r, binary.BigEndian, &b)
	return b
}

func readU64(r io.Reader) uint64 {
	var b uint64
	binary.Read(r, binary.BigEndian, &b)
	return b
}

const (
	TEST_COMPRESSION_FLAG CompressionFlag = 32
	TEST_ENCRYPTION_FLAG  EncryptionFlag  = 64
)

type testCompression struct {
}

// Flag implements CompressionHandler.
func (n *testCompression) Flag() CompressionFlag {
	return TEST_COMPRESSION_FLAG
}

func (n *testCompression) Decompress(b []byte) error {
	return nil
}

func (n *testCompression) Compress(b []byte) error {
	return nil
}

type testEncryptionFlag struct {
}

// Flag implements EncryptionHandler.
func (n *testEncryptionFlag) Flag() EncryptionFlag {
	return TEST_ENCRYPTION_FLAG
}

func (n *testEncryptionFlag) Decrypt(b []byte) error {
	return nil
}

func (n *testEncryptionFlag) Encrypt(b []byte) error {
	return nil
}

func Test_calculatePaddingLength(t *testing.T) {
	require.EqualValues(t, 0, calculatePaddingLength(0))
	require.EqualValues(t, 0, calculatePaddingLength(16))
	require.EqualValues(t, 8, calculatePaddingLength(8))
	require.EqualValues(t, 9, calculatePaddingLength(7))
}
