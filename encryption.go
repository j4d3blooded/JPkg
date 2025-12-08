package jpkg

type EncryptionFlag uint8

const (
	ENCRYPTION_NONE EncryptionFlag = iota
)

type EncryptionHandler interface {
	Flag() EncryptionFlag
	Decrypt([]byte) error
	Encrypt([]byte) error
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

func (n *NullEncryptionHandler) Encrypt(b []byte) error {
	return nil
}
