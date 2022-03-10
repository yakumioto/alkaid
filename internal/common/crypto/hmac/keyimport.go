/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package hmac

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type KeyImporter struct{}

func NewKey(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	return new(KeyImporter).KeyImport(raw, opts)
}

func (k *KeyImporter) KeyImport(raw interface{}, opts crypto.KeyImportOpts) (crypto.Key, error) {
	var key []byte

	switch raw := raw.(type) {
	case []byte:
		key = raw
	case string:
		key = []byte(raw)
	default:
		return nil, fmt.Errorf("only supports string or []byte type of key")
	}

	switch opts.Algorithm() {
	case crypto.HmacSha256, crypto.HmacSha512:
		return &Key{
			key:       key,
			algorithm: opts.Algorithm(),
		}, nil
	}

	return nil, fmt.Errorf("unsupported aes algorithm: %v", opts.Algorithm())
}
