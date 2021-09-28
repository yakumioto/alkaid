package aes

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func NewKeyImporter() *keyImporter {
	return &keyImporter{}
}

type keyImporter struct{}

func (kg *keyImporter) KeyImport(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	privateKey, ok := raw.([]byte)
	if !ok && privateKey != nil {
		return nil, fmt.Errorf("only supports []byte type of key")
	}

	switch opts.Algorithm() {
	case crypto.AES128, crypto.AES192, crypto.AES256:
		return &aesCBCPrivateKey{
			privateKey: privateKey,
			algorithm:  opts.Algorithm(),
		}, nil
	}

	return nil, fmt.Errorf("unsupported aes algorithm: %v", opts.Algorithm())
}
