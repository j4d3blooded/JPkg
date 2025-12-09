package jpkg_impl

import (
	"io"
)

func newNopWriterCloser(w io.Writer) io.WriteCloser {
	return &nopWriterCloser{w}
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
