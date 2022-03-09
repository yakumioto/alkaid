/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

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
	var privateKey []byte

	switch key := raw.(type) {
	case []byte:
		privateKey = key
	case string:
		privateKey = []byte(key)
	default:
		return nil, fmt.Errorf("only supports string or []byte type of key")
	}

	keyLen := 0

	switch opts.Algorithm() {
	case crypto.AES128:
		keyLen = 128 / 8
	case crypto.AES192:
		keyLen = 192 / 8
	case crypto.AES256:
		keyLen = 256 / 8
	}

	if len(privateKey) != keyLen {
		return nil, fmt.Errorf("the required key length is %v", keyLen)
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
