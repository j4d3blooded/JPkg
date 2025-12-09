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

type wrapReader struct {
	r io.ReadSeeker
}

// Read implements io.ReadSeeker.
func (w *wrapReader) Read(p []byte) (n int, err error) {
	return w.r.Read(p)
}

// Seek implements io.ReadSeeker.
func (w *wrapReader) Seek(offset int64, whence int) (int64, error) {
	return w.r.Seek(offset, whence)
}

func (w *wrapReader) readN(n uint64) []byte {
	chars := make([]byte, n)
	io.ReadFull(w, chars)
	return chars
}

func (w *wrapReader) u8() uint8 {
	return readT[uint8](w)
}

func (w *wrapReader) u64() uint64 {
	return readT[uint64](w)
}

func (w *wrapReader) readStr(n uint64) string {
	b := w.readN(n)
	return string(b)
}

func readT[T any](r io.Reader) T {
	t := new(T)
	binary.Read(r, binary.BigEndian, t)
	return *t
}

func calculatePaddingLength(M uint64) uint64 {
	padding := M / 16
	padding = padding + 1
	padding = 16 * padding
	padding = padding - M
	padding = padding % uint64(16)
	return padding
}
