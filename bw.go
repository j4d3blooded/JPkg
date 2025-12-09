package jpkg

import (
	"bytes"
	"encoding/binary"
	"io"
)

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
