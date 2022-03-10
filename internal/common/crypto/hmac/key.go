/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package hmac

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type Key struct {
	key       []byte
	algorithm crypto.Algorithm
}

func (h *Key) Bytes() ([]byte, error) {
	return nil, errors.New("not supported")
}

func (h *Key) SKI() []byte {
	sha := sha256.New()
	sha.Write(h.key)
	return sha.Sum(nil)
}

func (h *Key) Symmetric() bool {
	return false
}

func (h *Key) Private() bool {
	return true
}

func (h *Key) PublicKey() (crypto.Key, error) {
	return nil, errors.New("cannot call this method on a hmac key")
}

func (h *Key) Sign(digest []byte) ([]byte, error) {
	var hc hash.Hash
	switch h.algorithm {
	case crypto.HmacSha256:
		hc = hmac.New(sha256.New, h.key)
	case crypto.HmacSha512:
		hc = hmac.New(sha512.New, h.key)
	default:
		return nil, fmt.Errorf("not support %v algorithm", h.algorithm)
	}

	hc.Write(digest)

	return hc.Sum(nil), nil
}

func (h *Key) Verify(hash, sig []byte) bool {
	digest, err := h.Sign(hash)
	if err != nil {
		return false
	}

	return bytes.Equal(digest, sig)
}

func (h *Key) Encrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("cannot call this method on a hmac hash")
}

func (h *Key) Decrypt(_ []byte) ([]byte, error) {
	return nil, errors.New("cannot call this method on a hmac hash")
}
