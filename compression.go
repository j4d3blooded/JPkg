package jpkg

type CompressionFlag uint8

const (
	COMPRESSION_NONE CompressionFlag = iota
)

type CompressionHandler interface {
	Flag() CompressionFlag
	Decompress([]byte) error
	Compress([]byte) error
}

type NullCompressionHandler struct {
}

// Flag implements CompressionHandler.
func (n *NullCompressionHandler) Flag() CompressionFlag {
	return COMPRESSION_NONE
}

func (n *NullCompressionHandler) Decompress(b []byte) error {
	return nil
}

func (n *NullCompressionHandler) Compress(b []byte) error {
	return nil
}
