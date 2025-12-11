package jpkg_impl

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"io"
)

type CompressionFlag uint8

const (
	COMPRESSION_NONE CompressionFlag = iota
	COMPRESSION_LZW
)

type CompressionHandler interface {
	Flag() CompressionFlag
	Decompress(compressed []byte) ([]byte, error)
	Compress(uncompressed []byte) ([]byte, error)
}

type NullCompressionHandler struct {
}

// Flag implements CompressionHandler.
func (n *NullCompressionHandler) Flag() CompressionFlag {
	return COMPRESSION_NONE
}

func (n *NullCompressionHandler) Decompress(compressed []byte) ([]byte, error) {
	return compressed, nil
}

func (n *NullCompressionHandler) Compress(uncompressed []byte) ([]byte, error) {
	return uncompressed, nil
}

type LZWCompressionHandler struct {
}

// Flag implements CompressionHandler.
func (n *LZWCompressionHandler) Flag() CompressionFlag {
	return COMPRESSION_LZW
}

func (n *LZWCompressionHandler) Decompress(compressed []byte) ([]byte, error) {
	output := bytes.NewReader(compressed)
	reader := lzw.NewReader(output, lzw.LSB, 8)
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error during lzw decompression: %w", err)
	}
	if err := reader.Close(); err != nil {
		return nil, fmt.Errorf("error closing lzw decompression: %w", err)
	}
	return b, nil
}

func (n *LZWCompressionHandler) Compress(uncompressed []byte) ([]byte, error) {
	output := &bytes.Buffer{}
	writer := lzw.NewWriter(output, lzw.LSB, 8)
	_, err := writer.Write(uncompressed)
	if err != nil {
		return nil, fmt.Errorf("error during lzw compression: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing lzw compression: %w", err)
	}
	return output.Bytes(), nil
}
