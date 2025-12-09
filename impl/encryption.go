package jpkg_impl

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
)

type EncryptionFlag uint8

const (
	ENCRYPTION_NONE EncryptionFlag = iota
	ENCRYPTION_AES
)

type EncryptionHandler interface {
	Flag() EncryptionFlag
	Decrypt([]byte) error
	Encrypt(output io.Writer) (io.WriteCloser, error)
}

type NullEncryptionHandler struct {
}

// Flag implements EncryptionHandler.
func (n *NullEncryptionHandler) Flag() EncryptionFlag {
	return ENCRYPTION_NONE
}

func (n *NullEncryptionHandler) Decrypt(b []byte) error {
	return nil
}

func (n *NullEncryptionHandler) Encrypt(output io.Writer) (io.WriteCloser, error) {
	if wc, isWC := output.(io.WriteCloser); isWC {
		return wc, nil
	}
	return newNopWriterCloser(output), nil
}

type AESEncryptionHandler struct {
	Key []byte
}

type aesWriter struct {
	aead cipher.AEAD
	w    io.Writer
}

func (a *aesWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	b := a.aead.Seal(p[:0], nil, p, nil)
	_, err = a.w.Write(b)
	return
}

func (a *aesWriter) Close() error {
	return nil
}

// Flag implements EncryptionHandler.
func (n *AESEncryptionHandler) Flag() EncryptionFlag {
	return ENCRYPTION_AES
}

func (n *AESEncryptionHandler) Decrypt(b []byte) error {
	return nil
}

func (n *AESEncryptionHandler) Encrypt(output io.Writer) (io.WriteCloser, error) {
	block, err := aes.NewCipher(n.Key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %w", err)
	}
	aead, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher thingy: %w", err)
	}

	return &aesWriter{aead, output}, nil
}
