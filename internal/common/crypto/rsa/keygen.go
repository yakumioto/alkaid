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
	"fmt"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type keyGenerator struct{}

func KeyGen(opts crypto.KeyGenOpts) (crypto.Key, error) {
	return new(keyGenerator).KeyGen(opts)
}

func (k *keyGenerator) KeyGen(opts crypto.KeyGenOpts) (crypto.Key, error) {
	bits := 2048
	switch opts.Algorithm() {
	case crypto.Rsa1024:
		bits = 1024
	case crypto.Rsa2048:
		bits = 2048
	case crypto.Rsa4096:
		bits = 4096

	}

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("generating RSA key for [%v] error: [%s]", bits, err)
	}

	return &PrivateKey{privateKey: privateKey}, nil
}
