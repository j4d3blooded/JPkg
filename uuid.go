package jpkg

import (
	"crypto/rand"
	"io"
)

type UUID [16]byte

func NewUUID() UUID {
	var uuid UUID
	io.ReadFull(rand.Reader, uuid[:])
	return uuid
}
