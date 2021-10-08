/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func NewKeyGenerator() *keyGenerator {
	return &keyGenerator{}
}

type keyGenerator struct{}

func (kg *keyGenerator) KeyGen(opts crypto.KeyGenOpts) (crypto.Key, error) {
	var curve elliptic.Curve
	switch opts.Algorithm() {
	case crypto.ECDSAP256:
		curve = elliptic.P256()
	case crypto.ECDSAP384:
		curve = elliptic.P384()
	default:
		return nil, fmt.Errorf("unsupported ECDSA algorithm: %v", opts.Algorithm())
	}

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key for [%v] error: [%s]", curve, err)
	}

	return &ecdsaPrivateKey{privateKey}, nil
}
