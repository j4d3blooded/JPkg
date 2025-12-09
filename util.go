package jpkg

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type UUID [16]byte

func NewUUIDV4() UUID {
	var uuid UUID
	io.ReadFull(rand.Reader, uuid[:])
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return uuid
}

func serializeMetadataToJSON(data any) (string, uint64, error) {
	if data == nil {
		data = struct{}{}
	}

	metadataJsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", 0, fmt.Errorf("error parsing metadata as json: %w", err)
	}

	str := string(metadataJsonBytes)
	len := uint64(len(str))
	return str, len, nil
}

type nopWriterCloser struct {
	w io.Writer
}

// Close implements io.WriteCloser.
func (n *nopWriterCloser) Close() error {
	return nil
}

// Write implements io.WriteCloser.
func (np *nopWriterCloser) Write(p []byte) (n int, err error) {
	return np.w.Write(p)
}

func newNopWriterCloser(w io.Writer) io.WriteCloser {
	return &nopWriterCloser{w}
}

func newBufferWriter() *bufferWriter {
	return &bufferWriter{
		b: &bytes.Buffer{},
	}
}

type bufferWriter struct {
	b *bytes.Buffer
}

func (bw *bufferWriter) wb(data any) {
	binary.Write(bw.b, binary.BigEndian, data)
}

func (bw *bufferWriter) ws(data string) {
	io.WriteString(bw.b, data)
}

func (bw *bufferWriter) pad(count uint64) {
	db := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	for i := range int(count) {
		n := i % len(db)
		bw.wb(db[n])
	}
}
