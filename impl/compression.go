package jpkg_impl

import (
	"compress/lzw"
	"io"
)

type CompressionFlag uint8

const (
	COMPRESSION_NONE CompressionFlag = iota
	COMPRESSION_LZW
)

type CompressionHandler interface {
	Flag() CompressionFlag
	Decompress(r io.Reader) io.Reader
	Compress(output io.Writer) io.WriteCloser
}

type NullCompressionHandler struct {
}

// Flag implements CompressionHandler.
func (n *NullCompressionHandler) Flag() CompressionFlag {
	return COMPRESSION_NONE
}

func (n *NullCompressionHandler) Decompress(r io.Reader) io.Reader {
	return nil
}

func (n *NullCompressionHandler) Compress(output io.Writer) io.WriteCloser {
	if wc, isWC := output.(io.WriteCloser); isWC {
		return wc
	}
	return newNopWriterCloser(output)
}

type LZWCompressionHandler struct {
}

// Flag implements CompressionHandler.
func (n *LZWCompressionHandler) Flag() CompressionFlag {
	return COMPRESSION_LZW
}

func (n *LZWCompressionHandler) Decompress(r io.Reader) io.Reader {
	return nil
}

func (n *LZWCompressionHandler) Compress(output io.Writer) io.WriteCloser {
	return lzw.NewWriter(output, lzw.LSB, 8)
}
