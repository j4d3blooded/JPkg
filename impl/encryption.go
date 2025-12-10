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
	Decrypt(output io.Writer) (io.WriteCloser, error)
	Encrypt(output io.Writer) (io.WriteCloser, error)
}

type NullEncryptionHandler struct {
}

func (n *NullEncryptionHandler) Flag() EncryptionFlag {
	return ENCRYPTION_NONE
}

func (n *NullEncryptionHandler) Decrypt(output io.Writer) (io.WriteCloser, error) {
	return n.Encrypt(output)
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

type aesEncryptor struct {
	aead cipher.AEAD
	w    io.Writer
}

type aesDecryptor struct {
	aead cipher.AEAD
	w    io.Writer
}

func (a *aesEncryptor) Write(p []byte) (n int, err error) {
	n = len(p)
	b := a.aead.Seal(p[:0], nil, p, nil)
	_, err = a.w.Write(b)
	return
}

func (a *aesDecryptor) Write(p []byte) (n int, err error) {
	n = len(p)
	b, err := a.aead.Open(p[:0], nil, p, nil)
	if err != nil {
		err = fmt.Errorf("error during aes decryption: %w", err)
		return
	}
	_, err = a.w.Write(b)
	return
}

func (a *aesEncryptor) Close() error {
	return nil
}

func (n *AESEncryptionHandler) Flag() EncryptionFlag {
	return ENCRYPTION_AES
}

func (n *AESEncryptionHandler) Decrypt(output io.Writer) (io.WriteCloser, error) {
	return n.Encrypt(output)
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

	return &aesEncryptor{aead, output}, nil
}
