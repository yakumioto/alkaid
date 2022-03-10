/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package utils

import (
	"testing"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/hmac"
)

func TestEncrypt(t *testing.T) {
	aKey, _ := aes.NewKey("test password", &crypto.AES256KeyImportOpts{})
	hKey, _ := hmac.NewKey("test password", &crypto.HMACSha256ImportOpts{})

	ciphertext, _ := Encrypt(AesCbc256HmacSha256B64, []byte("hello word"), aKey, hKey)
	t.Log(ciphertext)
	data, _ := Decrypt(ciphertext, aKey, hKey)
	t.Log(string(data))
}
