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
	"github.com/yakumioto/alkaid/internal/common/factory"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/util"
	"github.com/yakumioto/alkaid/internal/errors"
)

var (
	logger = log.GetPackageLogger("services.users")
)

func InitializeRootUser(u *User) error {
	if err := u.create(); err != nil {
		logger.Errorf("[%v] initialize root user error: %v", u.ID, err)
		return errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to initialize root user")
	}

	return nil
}

type CreateRequest struct {
	ID                  string `json:"id,omitempty" validate:"required"`
	OrganizationID      string `json:"organizationId" validate:"required"`
	Name                string `json:"name" validate:"required"`
	Email               string `json:"email" validate:"required,email"`
	Password            string `json:"password" validate:"required"`
	TransactionPassword string `json:"transactionPassword" validate:"required"` // 交易密码仅用来加解密 PrivateKey
	Role                string `json:"role" validate:"required,oneof=user networkAdministrator organizationAdministrator"`
}

func Create(req *CreateRequest, userCtx *UserContext) (*User, error) {
	u, err := newUserByCreateRequest(req, userCtx)
	if err != nil {
		logger.Errorf("[%v] init create request error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusBadRequest, errors.ErrUserCreateVerifying,
			"init create request error")
	}

	sigPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		logger.Errorf("[%v] generate signature key error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate signature key")
	}
	sigPrivateKeyPem, err := sigPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the signature key to pem format error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the signature key to pem format")
	}

	tlsPrivateKey, err := factory.CryptoKeyGen(crypto.ECDSAP256)
	if err != nil {
		logger.Errorf("[%v] generate tls key error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate tls key")
	}
	tlsPrivateKeyPem, err := tlsPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the tls key to pem format error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the tls key to pem format")
	}

	aesKey, err := factory.CryptoKeyImport(req.TransactionPassword, crypto.AES256)
	if err != nil {
		logger.Errorf("[%v] import transaction password error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import transaction password key")
	}
	protectedSigPrivateKey, err := aesKey.Encrypt(sigPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption signing key error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption signing key failed")
	}
	protectedTLSPrivateKey, err := aesKey.Encrypt(tlsPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption tls key error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption tls key failed")
	}

	u.ProtectedSigPrivateKey = base64.StdEncoding.EncodeToString(protectedSigPrivateKey)
	u.ProtectedTLSPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)

	if err = u.create(); err != nil {
		logger.Errorf("[%v] create user error: %v", u.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to create user")
	}

	return u, nil
}

func GetDetailByID(id string) (*User, error) {
	user, err := FindByIDOrEmail(id)
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

func Login(req *LoginRequest) (*User, error) {
	logger.Debugf("entry ")
	user, err := FindByIDOrEmail(req.ID)
	if err != nil {
		if err == storage.ErrNotFound {
			logger.Infof("[%v] user not found", req.ID)
			return nil, errors.NewError(http.StatusNotFound, errors.ErrUserNotFount,
				"user not found")
		}

		logger.Errorf("[%v] query user error: %v", req.ID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"server unknown error")
	}

	if !util.ValidatePassword(req.Password, user.Email, user.Password) {
		logger.Infof("[%v] wrong user password", req.ID)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"wrong user password")
	}

	logger.Debugf("user: %v", user)
	return user, nil
}
