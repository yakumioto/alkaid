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

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type rsaPrivateKey struct {
	privateKey *rsa.PrivateKey
}

func (r *rsaPrivateKey) Bytes() ([]byte, error) {
	pkcs8Encoded := x509.MarshalPKCS1PrivateKey(r.privateKey)
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8Encoded}), nil
}

func (r *rsaPrivateKey) SKI() []byte {
	pubKey, _ := r.PublicKey()
	return pubKey.SKI()
}

func (r *rsaPrivateKey) Symmetric() bool {
	return false
}

func (r *rsaPrivateKey) Private() bool {
	return true
}

func (r *rsaPrivateKey) PublicKey() (crypto.Key, error) {
	return &rsaPublicKey{publicKey: &r.privateKey.PublicKey}, nil
}

func (r *rsaPrivateKey) Sign(digest []byte) ([]byte, error) {
	rsa.SignPSS(rand.Reader, r.privateKey, digest)
}

func (r *rsaPrivateKey) Verify(hash, sig []byte) bool {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPrivateKey) Encrypt(src []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPrivateKey) Decrypt(src []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privateKey, src, nil)
}

type rsaPublicKey struct {
	publicKey *rsa.PublicKey
}

func (r *rsaPublicKey) Bytes() ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) SKI() []byte {
	if r.publicKey == nil {
		return nil
	}

	// Marshal the public key and hash it
	raw := x509.MarshalPKCS1PublicKey(r.publicKey)
	hash := sha256.Sum256(raw)
	return hash[:]
}

func (r *rsaPublicKey) Symmetric() bool {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) Private() bool {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) PublicKey() (crypto.Key, error) {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) Sign(digest []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) Verify(hash, sig []byte) bool {
	// TODO implement me
	panic("implement me")
}

func (r *rsaPublicKey) Encrypt(src []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.publicKey, src, nil)
}

func (r *rsaPublicKey) Decrypt(src []byte) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}
