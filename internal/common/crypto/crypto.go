/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package crypto

import "crypto"

var (
	SHA256 = crypto.SHA256
)

type Key interface {
	Bytes() ([]byte, error)
	SKI() []byte
	Symmetric() bool
	Private() bool
	PublicKey() (Key, error)
	Sign(digest []byte) ([]byte, error)
	Verify(hash, sig []byte) bool
	Encrypt(src []byte) ([]byte, error)
	Decrypt(src []byte) ([]byte, error)
}

type KeyGenerator interface {
	KeyGen(opts KeyGenOpts) (Key, error)
}

type KeyGenOpts interface {
	Algorithm() Algorithm
}

type KeyImporter interface {
	KeyImport(raw interface{}, opts KeyImportOpts) (Key, error)
}

type KeyImportOpts interface {
	Algorithm() Algorithm
}
