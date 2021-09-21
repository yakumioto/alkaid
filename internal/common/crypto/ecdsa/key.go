package ecdsa

import (
	"crypto/ecdsa"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type ecdsaPrivateKey struct {
	privateKey *ecdsa.PrivateKey
}

func (e *ecdsaPrivateKey) Bytes() ([]byte, error) {
	panic("implement me")
}

func (e *ecdsaPrivateKey) SKI() []byte {
	panic("implement me")
}

func (e *ecdsaPrivateKey) Symmetric() bool {
	panic("implement me")
}

func (e *ecdsaPrivateKey) Private() bool {
	panic("implement me")
}

func (e *ecdsaPrivateKey) PublicKey() (crypto.Key, error) {
	panic("implement me")
}

type ecdsaPublicKey struct {
	privateKey *ecdsa.PublicKey
}

func (e *ecdsaPublicKey) Bytes() ([]byte, error) {
	panic("implement me")
}

func (e *ecdsaPublicKey) SKI() []byte {
	panic("implement me")
}

func (e *ecdsaPublicKey) Symmetric() bool {
	panic("implement me")
}

func (e *ecdsaPublicKey) Private() bool {
	panic("implement me")
}

func (e *ecdsaPublicKey) PublicKey() (crypto.Key, error) {
	panic("implement me")
}
