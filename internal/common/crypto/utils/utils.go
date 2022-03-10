/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/hmac"
)

type EncType int

const (
	AesCbc256B64 = iota
	AesCbc256HmacSha256B64
	Rsa2048OaepSha256B64
	Rsa2048OaepSha256HmacShaB64
)

func Encrypt(typ EncType, text []byte, keys ...interface{}) (string, error) {
	var (
		aesKey crypto.Key
		// rsaKey  crypto.Key
		hmacKey crypto.Key
	)

	for _, key := range keys {
		switch key := key.(type) {
		case *aes.CBCKey:
			aesKey = key
		case *hmac.Key:
			hmacKey = key
		}
	}

	switch typ {
	case AesCbc256HmacSha256B64:
		ciphertextBytes, err := aesKey.Encrypt(text)
		if err != nil {
			return "", err
		}
		ciphertext := base64.StdEncoding.EncodeToString(ciphertextBytes)
		sigBytes, err := hmacKey.Sign([]byte(ciphertext))
		if err != nil {
			return "", err
		}
		sig := base64.StdEncoding.EncodeToString(sigBytes)

		data := bytes.NewBuffer(nil)
		data.WriteString(strconv.Itoa(int(typ)))
		data.WriteString(".")
		data.WriteString(ciphertext)
		data.WriteString(".")
		data.WriteString(sig)

		return data.String(), nil
	}

	return "", nil
}

func Decrypt(text string, keys ...crypto.Key) ([]byte, error) {
	var (
		ak crypto.Key
		// rk  crypto.Key
		hk      crypto.Key
		sig     bool
		sigtext string
	)

	for _, key := range keys {
		switch key := key.(type) {
		case *aes.CBCKey:
			ak = key
		case *hmac.Key:
			hk = key
		}
	}

	texts := strings.Split(text, ".")
	if len(texts) < 2 && len(texts) > 3 {
		return nil, errors.New("irregular encrypted data format")
	}
	if len(texts) == 3 {
		sig = true
	}
	typ, err := strconv.Atoi(texts[0])
	if err != nil {
		return nil, errors.New("irregular encrypted data format")
	}
	ciphertext := texts[1]
	if sig {
		sigtext = texts[2]
	}

	switch EncType(typ) {
	case AesCbc256HmacSha256B64:
		// step1: hmac
		// step2: aes decrypt
		sigBytes, err := base64.StdEncoding.DecodeString(sigtext)
		if err != nil {
			return nil, fmt.Errorf("base64 decode sig error: %v", err)
		}
		if !hk.Verify([]byte(ciphertext), sigBytes) {
			return nil, fmt.Errorf("hmac verify error: %v", err)
		}

		ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return nil, fmt.Errorf("base64 decode ciphertext error: %v", err)
		}
		dataBytes, err := ak.Decrypt(ciphertextBytes)
		if err != nil {
			return nil, err
		}

		return dataBytes, nil
	}

	return nil, nil
}
