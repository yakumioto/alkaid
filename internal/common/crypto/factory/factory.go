/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package factory

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/ecdsa"
)

func CryptoKeyGen(algorithm crypto.Algorithm) (crypto.Key, error) {
	switch algorithm {
	case crypto.EcdsaP256:
		return ecdsa.KeyGen(&crypto.ECDSAP256KeyGenOpts{})
	case crypto.EcdsaP384:
		return ecdsa.KeyGen(&crypto.ECDSAP384KeyGenOpts{})
	}

	return nil, fmt.Errorf("not found key generator: %v", algorithm)
}

func CryptoKeyImport(raw interface{}, algorithm crypto.Algorithm) (crypto.Key, error) {
	switch algorithm {
	case crypto.AesCbc128:
		return aes.NewKey(raw, &crypto.AES128KeyImportOpts{})
	case crypto.AesCbc192:
		return aes.NewKey(raw, &crypto.AES192KeyImportOpts{})
	case crypto.AesCbc256:
		return aes.NewKey(raw, &crypto.AES256KeyImportOpts{})
	case crypto.EcdsaP256, crypto.EcdsaP384:
		return ecdsa.KeyImport(raw)
	}

	return nil, fmt.Errorf("not found key importer: %v", algorithm)
}
