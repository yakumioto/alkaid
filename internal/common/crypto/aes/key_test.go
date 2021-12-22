/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package aes

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func TestAESCBCPrivateEncryptAndDecrypt(t *testing.T) {
	tcs := []struct {
		password         string
		algorithm        crypto.Algorithm
		randomIVFunc     randomIVFunc
		text             string
		base64Ciphertext string
	}{
		{
			"aes128",
			crypto.AES128,
			func(len int) ([]byte, error) {
				return []byte{68, 161, 247, 38, 29, 93, 40, 189, 245, 37, 160, 123, 196, 137, 53, 69}, nil
			},
			"hello world",
			"AAAAAAAAAAAAAESh9yYdXSi99SWge8SJNUUhyXMNUZhrYOy4VLgs2BbB",
		},
		{
			"aes192",
			crypto.AES192,
			func(len int) ([]byte, error) {
				return []byte{116, 83, 64, 105, 45, 5, 144, 131, 16, 66, 7, 54, 175, 180, 140, 217}, nil
			},
			"hello world",
			"AAAAAAAAAAAAAHRTQGktBZCDEEIHNq+0jNl2LbrzzzKpmud2HXpgqAqX",
		},
		{
			"aes256",
			crypto.AES256,
			func(len int) ([]byte, error) {
				return []byte{89, 29, 100, 99, 47, 20, 209, 119, 122, 133, 110, 49, 231, 200, 124, 49}, nil
			},
			"hello world",
			"AAAAAAAAAAAAAFkdZGMvFNF3eoVuMefIfDFNsD3sPn5A49AsUnyyD2vG",
		},
	}

	for _, tc := range tcs {
		ack := &aesCBCPrivateKey{privateKey: []byte(tc.password), algorithm: tc.algorithm}
		randomIV = tc.randomIVFunc
		ciphertext, err := ack.Encrypt([]byte(tc.text))
		assert.NoError(t, err)
		assert.Equal(t, base64.StdEncoding.EncodeToString(ciphertext), tc.base64Ciphertext, "ciphertext error")
		text, err := ack.Decrypt(ciphertext)
		assert.NoError(t, err)
		assert.Equal(t, string(text), tc.text, "text error")
		t.Logf("text:%v; ciphertext: %v", string(text), base64.StdEncoding.EncodeToString(ciphertext))
	}
}
