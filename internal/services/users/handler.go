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
	"fmt"
	"net/http"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/factory"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/errors"
)

type CreateRequest struct {
	ID                  string `json:"id,omitempty" validate:"required"`
	OrganizationID      string `json:"organizationId" validate:"required"`
	Name                string `json:"name" validate:"required"`
	Email               string `json:"email" validate:"required,email"`
	Password            string `json:"password" validate:"required"`
	TransactionPassword string `json:"transactionPassword" validate:"required"` // 交易密码仅用来加解密 PrivateKey
	Role                string `json:"role" validate:"required,oneof=user networkAdministrator organizationAdministrator"`
}

func (u *User) Create(req *CreateRequest, userCtx *UserContext) error {
	u.initByCreateRequest(req, userCtx)

	sigPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to generate signature key", err)
	}
	sigPrivateKeyPem, err := sigPrivateKey.Bytes()
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to convert the signature key to pem format", err)
	}

	tlsPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to generate tls key", err)
	}
	tlsPrivateKeyPem, err := tlsPrivateKey.Bytes()
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to convert the tls key to pem format", err)
	}

	aesKey, err := factory.CryptoKeyImport(req.TransactionPassword, crypto.AES256)
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to import encryption key", err)
	}
	protectedSigPrivateKey, err := aesKey.Encrypt(sigPrivateKeyPem)
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"encryption signing key failed", err)
	}
	protectedTLSPrivateKey, err := aesKey.Encrypt(tlsPrivateKeyPem)
	if err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"encryption tls key failed", err)
	}

	u.ProtectedSigPrivateKey = base64.StdEncoding.EncodeToString(protectedSigPrivateKey)
	u.ProtectedTLSPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)

	if err = u.create(); err != nil {
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			"failed to create user", err)
	}

	return nil
}

func (u *User) GetDetailByID(id string) error {
	u.initUserByID(id)

	if err := u.findByID(); err != nil {
		if err == storage.ErrNotFound {
			return u.error(http.StatusNotFound, errors.UserNotFount,
				fmt.Sprintf("user [%v] not found", id), err)
		}
		return u.error(http.StatusInternalServerError, errors.ServerUnknownError,
			fmt.Sprintf("failed to query user [%v]", id), err)
	}

	return nil
}

func (u *User) error(statusCode int, code errors.Code, msg string, err error) error {
	log.Errorf("%s: %v", msg, err)
	return errors.NewError(statusCode, code, msg)
}
