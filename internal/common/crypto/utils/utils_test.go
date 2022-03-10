/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package utils

import (
	"github.com/yakumioto/alkaid/internal/common/crypto/rsa"
	"testing"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/hmac"
)

func TestEncrypt(t *testing.T) {
	rPrivKey, _ := rsa.KeyGen(&crypto.RSA2048KeyImportOpts{})
	rPubKey, _ := rPrivKey.PublicKey()
	aKey, _ := aes.NewKey("test password", &crypto.AES256KeyImportOpts{})
	hKey, _ := hmac.NewKey("test password", &crypto.HMACSha256ImportOpts{})

	ciphertext, _ := Encrypt(AesCbc256HmacSha256B64, []byte("hello word"), aKey, hKey)
	t.Log(ciphertext)
	data, _ := Decrypt(ciphertext, aKey, hKey)
	t.Log(string(data))

	ciphertext, err := Encrypt(Rsa2048OaepSha256HmacShaB64, []byte("hello word"), rPubKey, hKey)
	if err != nil {
		t.Log(err)
	}
	t.Log(ciphertext)
	data, _ = Decrypt(ciphertext, rPrivKey, hKey)
	t.Log(string(data))
}
