package jpkg

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
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

const MAGIC_NUMBER = uint32(0x6A706B67)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func normalizeFilePath(name string) string {
	name = strings.TrimPrefix(name, ".")

	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	name = filepath.Clean(name)
	name = strings.ReplaceAll(name, "/", "\\")
	return name
}
