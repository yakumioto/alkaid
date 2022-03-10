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
	"net/http"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/factory"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/utils"
	"github.com/yakumioto/alkaid/internal/errors"
)

var (
	logger = log.GetPackageLogger("services.users")
)

type CreateRequest struct {
	UserID              string `json:"userId" validate:"required"`
	Name                string `json:"name" validate:"required"`
	Email               string `json:"email" validate:"required,email"`
	Root                bool   `json:"root"`
	Password            string `json:"password" validate:"required"`
	TransactionPassword string `json:"transactionPassword" validate:"required"` // 交易密码仅用来加解密 PrivateKey
}

func Create(req *CreateRequest) (*User, error) {
	u := newUserByCreateRequest(req)

	signPrivateKey, err := factory.CryptoKeyGen(crypto.EcdsaP256)
	if err != nil {
		logger.Errorf("[%v] generate signature key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate signature key")
	}
	sigPrivateKeyPem, err := signPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the signature key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the signature key to pem format")
	}

	tlsPrivateKey, err := factory.CryptoKeyGen(crypto.EcdsaP256)
	if err != nil {
		logger.Errorf("[%v] generate tls key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate tls key")
	}
	tlsPrivateKeyPem, err := tlsPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the tls key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the tls key to pem format")
	}

	aesKey, err := factory.CryptoKeyImport([]byte(req.TransactionPassword), crypto.AesCbc256)
	if err != nil {
		logger.Errorf("[%v] import transaction password error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import transaction password key")
	}
	protectedSigPrivateKey, err := aesKey.Encrypt(sigPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption signing key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption signing key failed")
	}
	protectedTLSPrivateKey, err := aesKey.Encrypt(tlsPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption tls key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption tls key failed")
	}

	u.ProtectedSignPrivateKey = base64.StdEncoding.EncodeToString(protectedSigPrivateKey)
	u.ProtectedTLSPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)

	if err = u.Create(); err != nil {
		logger.Errorf("[%v] create user error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to create user")
	}

	return u, nil
}

func GetDetailByID(id string) (*User, error) {
	user, err := FindUserByID(id)
	if err != nil {
		if err == storage.ErrNotFound {
			logger.Warnf("[%v] user not found", id)
			return nil, errors.NewError(http.StatusNotFound, errors.ErrUserNotFount,
				"user not found")
		}
		logger.Errorf("[%v] query user error: %v", id, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"server unknown error")
	}

	return user, nil
}

type LoginRequest struct {
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

func Login(req *LoginRequest) (*User, []*UserOrganizations, error) {
	user, err := FindUserByID(req.ID)
	if err != nil {
		if err == storage.ErrNotFound {
			logger.Infof("[%v] user not found", req.ID)
			return nil, nil, errors.NewError(http.StatusNotFound, errors.ErrUserNotFount,
				"user not found")
		}

		logger.Errorf("[%v] query user error: %v", req.ID, err)
		return nil, nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"server unknown error")
	}

	if !utils.ValidatePassword(req.Password, user.Email, user.Password) {
		logger.Infof("[%v] wrong user password", req.ID)
		return nil, nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"wrong user password")
	}

	organizations, err := FindUserOrganizationsByUserID(user.UserID)
	if err != nil && err != storage.ErrNotFound {
		logger.Errorf("[%v] query user organizations error: %v", req.ID, err)
		return nil, nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"server unknown error")
	}

	return user, organizations, nil
}
