package jpkg

import (
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

func serializeMetadataToJSON(data any) (string, error) {
	if data == nil {
		data = struct{}{}
	}

	metadataJsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error parsing metadata as json: %w", err)
	}

	str := string(metadataJsonBytes)
	return str, nil
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

const MAGIC_NUMBER = uint32(0x6A706B67)
