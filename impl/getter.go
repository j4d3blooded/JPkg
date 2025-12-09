package jpkg_impl

import "fmt"

func GetCompressionHandler(flag CompressionFlag) CompressionHandler {
	switch flag {
	case COMPRESSION_NONE:
		return &NullCompressionHandler{}
	case COMPRESSION_LZW:
		return &LZWCompressionHandler{}
	}

	panic(fmt.Errorf("invalid compression flag: %v", flag))
}

func GetEncryptionHandler(flag EncryptionFlag, key []byte) EncryptionHandler {
	switch flag {
	case ENCRYPTION_NONE:
		return &NullEncryptionHandler{}
	case ENCRYPTION_AES:
		return &AESEncryptionHandler{key}
	}

	panic(fmt.Errorf("invalid encryption flag: %v", flag))
}

func GetCryptoHandler(flag CryptoFlag) CryptoHandler {
	switch flag {
	case CRYPTO_NONE:
		return &NullCryptoHandler{}
	}

	panic(fmt.Errorf("invalid crypto flag: %v", flag))
}

func GetHashHandler(flag HasherFlag) HasherHandler {
	switch flag {
	case HASHER_NONE:
		return &NullHasherHandler{}
	}

	panic(fmt.Errorf("invalid hash flag: %v", flag))
}
