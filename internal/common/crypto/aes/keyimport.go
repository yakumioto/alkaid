/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package aes

import (
	"crypto/sha256"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"golang.org/x/crypto/pbkdf2"
)

func NewKey(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	return new(keyImporter).KeyImport(raw, opts)
}

type keyImporter struct{}

func (kg *keyImporter) KeyImport(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	var key []byte

	switch raw := raw.(type) {
	case []byte:
		key = raw
	case string:
		key = []byte(raw)
	default:
		return nil, fmt.Errorf("only supports string or []byte type of key")
	}

	keyLen := 0

	switch opts.Algorithm() {
	case crypto.AesCbc128:
		keyLen = 128 / 8
	case crypto.AesCbc192:
		keyLen = 192 / 8
	case crypto.AesCbc256:
		keyLen = 256 / 8
	}

	if len(key) != keyLen {
		key = pbkdf2.Key(key, key, 1000, keyLen, sha256.New)
	}

	switch opts.Algorithm() {
	case crypto.AesCbc128, crypto.AesCbc192, crypto.AesCbc256:
		return &CBCKey{
			key: key,
		}, nil
	}

	return nil, fmt.Errorf("unsupported aes algorithm: %v", opts.Algorithm())
}
