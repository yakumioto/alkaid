package factory

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/ecdsa"
)

func CryptoKeyGen(algorithm crypto.Algorithm) (crypto.Key, error) {
	switch algorithm {
	case crypto.ECDSAP256:
		return ecdsa.NewKeyGenerator().KeyGen(&crypto.ECDSAP256KeyGenOpts{})
	case crypto.ECDSAP384:
		return ecdsa.NewKeyGenerator().KeyGen(&crypto.ECDSAP384KeyGenOpts{})
	}

	return nil, fmt.Errorf("not found key generator: %v", algorithm)
}

func CryptoKeyImport(raw interface{}, algorithm crypto.Algorithm) (crypto.Key, error) {
	switch algorithm {
	case crypto.AES128:
		return aes.NewKeyImporter().KeyImport(raw, &crypto.AES128KeyImportOpts{})
	case crypto.AES192:
		return aes.NewKeyImporter().KeyImport(raw, &crypto.AES192KeyImportOpts{})
	case crypto.AES256:
		return aes.NewKeyImporter().KeyImport(raw, &crypto.AES256KeyImportOpts{})
	}

	return nil, fmt.Errorf("not found key importer: %v", algorithm)
}
