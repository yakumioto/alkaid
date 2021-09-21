package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type ecdsaKeyGenerator struct{}

func (kg *ecdsaKeyGenerator) KeyGen(opts crypto.KeyGenOpts) (crypto.Key, error) {
	var curve elliptic.Curve
	switch opts.Algorithm() {
	case P256:
		curve = elliptic.P256()
	case P384:
		curve = elliptic.P384()
	default:
		return nil, fmt.Errorf("unsupported ecdsa algorithm: %v", opts.Algorithm())
	}

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key for [%v] error: [%s]", curve, err)
	}

	return &ecdsaPrivateKey{privateKey}, nil
}
