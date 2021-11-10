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

func TestAesCBCPrivateKey_Encrypt(t *testing.T) {
	ack := &aesCBCPrivateKey{privateKey: []byte("test key"), algorithm: crypto.AES256}
	ciphertext, err := ack.Encrypt([]byte("hello world"))
	assert.NoError(t, err)
	t.Log(base64.StdEncoding.EncodeToString(ciphertext))
}

func TestAesCBCPrivateKey_Decrypt(t *testing.T) {
	ciphertext, _ := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAOqju1sRnH0oOMdiJuO2XKY7UReOuV42/bpJcYXohLhf")
	ack := &aesCBCPrivateKey{privateKey: []byte("test key"), algorithm: crypto.AES256}
	text, err := ack.Decrypt(ciphertext)
	assert.NoError(t, err)
	t.Log(string(text))
}
