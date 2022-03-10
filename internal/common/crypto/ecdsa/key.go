/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package ecdsa

import (
	cryptoStd "crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type ecdsaPrivateKey struct {
	privateKey *ecdsa.PrivateKey
}

func (e *ecdsaPrivateKey) Bytes() ([]byte, error) {
	pkcs8Encoded, err := x509.MarshalPKCS8PrivateKey(e.privateKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal private key")
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8Encoded}), nil
}

func (e *ecdsaPrivateKey) SKI() []byte {
	pubKey, _ := e.PublicKey()
	return pubKey.SKI()
}

func (e *ecdsaPrivateKey) Symmetric() bool {
	return false
}

func (e *ecdsaPrivateKey) Private() bool {
	return true
}

func (e *ecdsaPrivateKey) PublicKey() (crypto.Key, error) {
	return &ecdsaPublicKey{publicKey: &e.privateKey.PublicKey}, nil
}

func (e *ecdsaPrivateKey) Sign(digest []byte) ([]byte, error) {
	return e.privateKey.Sign(rand.Reader, digest, cryptoStd.SHA256)
}

func (e *ecdsaPrivateKey) Verify(hash, sig []byte) bool {
	return ecdsa.VerifyASN1(&e.privateKey.PublicKey, hash, sig)
}

func (e *ecdsaPrivateKey) Encrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

func (e *ecdsaPrivateKey) Decrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

type ecdsaPublicKey struct {
	publicKey *ecdsa.PublicKey
}

func (e *ecdsaPublicKey) Bytes() ([]byte, error) {
	pkcs8Encoded, err := x509.MarshalPKIXPublicKey(e.publicKey)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal public key")
	}

	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pkcs8Encoded}), nil
}

func (e *ecdsaPublicKey) SKI() []byte {
	// Marshall the public key
	raw := elliptic.Marshal(e.publicKey.Curve, e.publicKey.X, e.publicKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (e *ecdsaPublicKey) Symmetric() bool {
	return false
}

func (e *ecdsaPublicKey) Private() bool {
	return false
}

func (e *ecdsaPublicKey) PublicKey() (crypto.Key, error) {
	return e, nil
}

func (e *ecdsaPublicKey) Sign(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

func (e *ecdsaPublicKey) Verify(hash, sig []byte) bool {
	return ecdsa.VerifyASN1(e.publicKey, hash, sig)
}

func (e *ecdsaPublicKey) Encrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

func (e *ecdsaPublicKey) Decrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}
