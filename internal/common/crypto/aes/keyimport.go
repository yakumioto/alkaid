package aes

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type KeyImport struct{}

func (kg *KeyImport) KeyImport(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	privateKey, ok := raw.([]byte)
	if !ok && privateKey != nil {
		return nil, fmt.Errorf("only supports []byte type of key")
	}

	switch opts.Algorithm() {
	case crypto.AES128:
	case crypto.AES192:
	case crypto.AES256:
	default:
		return nil, fmt.Errorf("unsupported aes algorithm: %v", opts.Algorithm())
	}

	return &aesCBCPrivateKey{
		privateKey: privateKey,
		algorithm:  opts.Algorithm(),
	}, nil
}
