/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package users

import (
	"encoding/base64"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/factory"
)

type CreateRequest struct {
	ID                  string `json:"id,omitempty" validate:"required"`
	Email               string `json:"email,omitempty" validate:"required,email"`
	Password            string `json:"password,omitempty" validate:"required"`
	TransactionPassword string `json:"transactionPassword,omitempty" validate:"required"` // 交易密码仅用来加解密 PrivateKey
	Role                string `json:"role,omitempty" validate:"required"`
}

func CreateUser(req *CreateRequest) (*User, error) {
	user := newUser(req)

	sigPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		return nil, err
	}
	sigPrivateKeyPem, err := sigPrivateKey.Bytes()
	if err != nil {
		return nil, err
	}
	tlsPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		return nil, err
	}
	tlsPrivateKeyPem, err := tlsPrivateKey.Bytes()
	if err != nil {
		return nil, err
	}

	aesKey, err := factory.CryptoKeyImport(req.TransactionPassword, crypto.AES256)
	if err != nil {
		return nil, err
	}
	protectedSigPrivateKey, err := aesKey.Encrypt(sigPrivateKeyPem)
	if err != nil {
		return nil, err
	}
	protectedTLSPrivateKey, err := aesKey.Encrypt(tlsPrivateKeyPem)
	if err != nil {
		return nil, err
	}

	user.ProtectedSigPrivateKey = base64.StdEncoding.EncodeToString(protectedSigPrivateKey)
	user.ProtectedTLSPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)

	if err := user.create(); err != nil {
		return nil, err
	}

	return user, nil
}
