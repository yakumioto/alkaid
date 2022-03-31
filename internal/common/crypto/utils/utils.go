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
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lithammer/shortuuid"
	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/aes"
	"github.com/yakumioto/alkaid/internal/common/crypto/hmac"
	"github.com/yakumioto/alkaid/internal/common/crypto/rsa"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

type EncType int

const (
	AesCbc256B64 = iota
	AesCbc256HmacSha256B64
	Rsa2048OaepSha256B64
	Rsa2048OaepSha256HmacShaB64
)

type StretchedKey struct {
	Enc []byte
	Mac []byte
}

func (s *StretchedKey) Key() []byte {
	return append(s.Enc, s.Mac...)
}

func GetStretchedKey(masterKey []byte) (*StretchedKey, error) {
	encData := make([]byte, 32)
	macData := make([]byte, 32)
	_, err := hkdf.Expand(sha256.New, masterKey, []byte("enc")).Read(encData)
	if err != nil {
		return nil, err
	}
	_, err = hkdf.Expand(sha256.New, masterKey, []byte("mac")).Read(macData)
	if err != nil {
		return nil, err
	}
	return &StretchedKey{
		Enc: encData,
		Mac: macData,
	}, nil
}

func GenSymmetricKey() (*StretchedKey, error) {
	encData := make([]byte, 32)
	macData := make([]byte, 32)
	_, err := rand.Read(encData)
	if err != nil {
		return nil, err
	}
	_, err = rand.Read(macData)
	if err != nil {
		return nil, err
	}

	return &StretchedKey{
		Enc: encData,
		Mac: macData,
	}, nil
}

func GetMasterKey(password, salt string) []byte {
	return pbkdf2.Key([]byte(password), []byte(salt), 100000, 32, sha256.New)
}

func HashPassword(password, salt string, iter int) string {
	passwordHash := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	passwordHashBase64 := base64.StdEncoding.EncodeToString(passwordHash)
	return fmt.Sprintf("%d.%s", iter, passwordHashBase64)
}

// PBKDF2WithSha256 用于生成扩展密钥
func PBKDF2WithSha256(password, salt []byte, keyLen int) []byte {
	return pbkdf2.Key(password, salt, 1, keyLen, sha256.New)
}

func ValidatePassword(password, salt, passwordHash string) bool {
	iter, err := strconv.Atoi(strings.SplitN(passwordHash, ".", 2)[0])
	if err != nil {
		return false
	}

	if passwordHash != HashPassword(password, salt, iter) {
		return false
	}

	return true
}

func GenResourceID(namespace string) string {
	return fmt.Sprintf("%s-%s", namespace, shortuuid.New())
}

func Encrypt(typ EncType, text []byte, keys ...interface{}) (string, error) {
	var (
		ak crypto.Key
		rk crypto.Key
		hk crypto.Key
	)

	for _, key := range keys {
		switch key := key.(type) {
		case *aes.CBCKey:
			ak = key
		case *hmac.Key:
			hk = key
		case *rsa.PublicKey:
			rk = key
		}
	}

	var (
		ciphertextBytes []byte
		err             error
	)
	switch typ {
	case AesCbc256HmacSha256B64:
		if ak == nil || hk == nil {
			return "", errors.New("not found aes key or hmac key")
		}
		ciphertextBytes, err = ak.Encrypt(text)
		if err != nil {
			return "", err
		}

	case Rsa2048OaepSha256HmacShaB64:
		if rk == nil || hk == nil {
			return "", errors.New("not found rsa key or hmac key")
		}
		ciphertextBytes, err = rk.Encrypt(text)
		if err != nil {
			return "", err
		}
	}

	ciphertext := base64.StdEncoding.EncodeToString(ciphertextBytes)
	sigBytes, err := hk.Sign([]byte(ciphertext))
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

func Decrypt(text string, keys ...crypto.Key) ([]byte, error) {
	var (
		ak      crypto.Key
		rk      crypto.Key
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
		case *rsa.PrivateKey:
			rk = key
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
		if hk == nil || ak == nil {
			return nil, errors.New("not found aes key or hmac key")
		}

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
	case Rsa2048OaepSha256HmacShaB64:
		if hk == nil || rk == nil {
			return nil, errors.New("not found aes key or hmac key")
		}
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
		dataBytes, err := rk.Decrypt(ciphertextBytes)
		if err != nil {
			return nil, err
		}

		return dataBytes, nil
	}

	return nil, nil
}
