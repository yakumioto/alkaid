/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type PrivateKey struct {
	privateKey *rsa.PrivateKey
}

func (r *PrivateKey) Bytes() ([]byte, error) {
	pkcs1Encoded := x509.MarshalPKCS1PrivateKey(r.privateKey)
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs1Encoded}), nil
}

func (r *PrivateKey) SKI() []byte {
	pubKey, _ := r.PublicKey()
	return pubKey.SKI()
}

func (r *PrivateKey) Symmetric() bool {
	return false
}

func (r *PrivateKey) Private() bool {
	return true
}

func (r *PrivateKey) PublicKey() (crypto.Key, error) {
	return &PublicKey{publicKey: &r.privateKey.PublicKey}, nil
}

func (r *PrivateKey) Sign(digest []byte) ([]byte, error) {
	return rsa.SignPSS(rand.Reader, r.privateKey, crypto.SHA256, digest, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
}

func (r *PrivateKey) Verify(_, _ []byte) bool {
	return false
}

func (r *PrivateKey) Encrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

func (r *PrivateKey) Decrypt(src []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, src, nil)
}

type PublicKey struct {
	publicKey *rsa.PublicKey
}

func (r *PublicKey) Bytes() ([]byte, error) {
	pkcs1Encoded := x509.MarshalPKCS1PublicKey(r.publicKey)
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pkcs1Encoded}), nil
}

func (r *PublicKey) SKI() []byte {
	if r.publicKey == nil {
		return nil
	}

	// Marshal the public key and hash it
	raw := x509.MarshalPKCS1PublicKey(r.publicKey)
	hash := sha256.Sum256(raw)
	return hash[:]
}

func (r *PublicKey) Symmetric() bool {
	return false
}

func (r *PublicKey) Private() bool {
	return false
}

func (r *PublicKey) PublicKey() (crypto.Key, error) {
	return r, nil
}

func (r *PublicKey) Sign(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}

func (r *PublicKey) Verify(hash, sig []byte) bool {
	if err := rsa.VerifyPSS(r.publicKey, crypto.SHA256, hash, sig, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	}); err != nil {
		return false
	}

	return true
}

func (r *PublicKey) Encrypt(src []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, src, nil)
}

func (r *PublicKey) Decrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("not supported")
}
