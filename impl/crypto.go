package jpkg_impl

type CryptoFlag uint8

const (
	CRYPTO_NONE CryptoFlag = iota
)

type CryptoHandler interface {
	Flag() CryptoFlag
}

type NullCryptoHandler struct {
}

func (n *NullCryptoHandler) Flag() CryptoFlag {
	return CRYPTO_NONE
}
