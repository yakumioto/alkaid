/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package rsa

import (
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type KeyImporter struct{}

func KeyImport(raw interface{}) (crypto.Key, error) {
	return new(KeyImporter).KeyImport(raw, nil)
}

func (k *KeyImporter) KeyImport(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	var der []byte

	switch raw := raw.(type) {
	case []byte:
		der = raw
	case string:
		der = []byte(raw)
	default:
		return nil, fmt.Errorf("only supports string or []byte type of key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(der)
	if err == nil {
		return &PrivateKey{privateKey: privKey}, nil
	}

	pubKey, err := x509.ParsePKCS1PublicKey(der)
	if err == nil {
		return &PublicKey{publicKey: pubKey}, nil
	}

	return nil, errors.New("is not ecdsa key")
}
