/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

type randomIVFunc func(len int) ([]byte, error)

var (
	randomIV randomIVFunc = func(len int) ([]byte, error) {
		iv := make([]byte, len)
		if _, err := rand.Read(iv); err != nil {
			return nil, err
		}

		return iv, nil
	}
)

type CBCKey struct {
	key []byte
}

func (a *CBCKey) Bytes() ([]byte, error) {
	return nil, errors.New("not supported")
}

func (a *CBCKey) SKI() []byte {
	hash := sha256.New()
	// hash.Write([]byte{0x01})
	hash.Write(a.key)
	return hash.Sum(nil)
}

func (a *CBCKey) Symmetric() bool {
	return true
}

func (a *CBCKey) Private() bool {
	return true
}

func (a *CBCKey) PublicKey() (crypto.Key, error) {
	return nil, errors.New("cannot call this method on a symmetric key")
}

func (a *CBCKey) Sign(_ []byte) ([]byte, error) {
	return nil, errors.New("cannot call this method on a symmetric key")
}

func (a *CBCKey) Verify(_, _ []byte) bool {
	return false
}

func (a *CBCKey) Encrypt(text []byte) ([]byte, error) {
	paddedText := pkcs7Padding(text)

	iv, err := randomIV(aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("random iv error: %v", err)
	}

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("new chipher error: %v", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(paddedText))
	mode.CryptBlocks(dst, paddedText)

	ciphertext := bytes.NewBuffer(nil)
	ciphertext.Write(iv)
	ciphertext.Write(dst)

	return ciphertext.Bytes(), nil
}

func (a *CBCKey) Decrypt(ciphertext []byte) ([]byte, error) {
	iv := ciphertext[0:16]
	src := ciphertext[16:]

	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, fmt.Errorf("new chipher error: %v", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	paddedText := make([]byte, len(src))
	mode.CryptBlocks(paddedText, src)

	return pkcs7UnPadding(paddedText), nil
}

func pkcs7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize

	var paddingText []byte
	if padding == 0 {
		paddingText = bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)
	} else {
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(src, paddingText...)
}

func pkcs7UnPadding(src []byte) []byte {
	unPadding := int(src[len(src)-1])
	return src[:(len(src) - unPadding)]
}
