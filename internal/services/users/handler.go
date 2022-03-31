/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package users

import (
	"net/http"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/factory"
	"github.com/yakumioto/alkaid/internal/common/crypto/utils"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/errors"
)

var (
	logger = log.GetPackageLogger("services.users")
)

type CreateRequest struct {
	ID       string `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Root     bool   `json:"root"`
	Password string `json:"password" validate:"required"`
}

func Create(req *CreateRequest) (*User, error) {
	u := newUserByCreateRequest(req)

	// 对用户生成扩展密钥
	stretchedKey, err := u.stretchedKey(req.Password)
	if err != nil {
		logger.Errorf("[%v] generate stretch key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate stretch key")
	}

	// 生成用户唯一的对称密钥
	symmetricKey, err := utils.GenSymmetricKey()
	if err != nil {
		logger.Errorf("[%v] generate symmetric key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate symmetric key")
	}

	signPrivateKey, err := factory.CryptoKeyGen(crypto.EcdsaP256)
	if err != nil {
		logger.Errorf("[%v] generate signature key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate signature key")
	}
	signPrivateKeyPem, err := signPrivateKey.Bytes()
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

	rsaPrivateKey, err := factory.CryptoKeyGen(crypto.Rsa2048)
	if err != nil {
		logger.Errorf("[%v] generate rsa key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate rsa key")
	}
	rsaPrivateKeyPem, err := rsaPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the rsa key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the rsa key to pem format")
	}

	symmetricKeyAesKey, err := factory.CryptoKeyImport(symmetricKey.Enc, crypto.AesCbc256)
	if err != nil {
		logger.Errorf("[%v] import aes password error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import aes password key")
	}
	symmetricKeyHashKey, err := factory.CryptoKeyImport(symmetricKey.Mac, crypto.HmacSha256)
	if err != nil {
		logger.Errorf("[%v] import hash password error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import hash password key")
	}

	protectedSignPrivateKey, err := utils.Encrypt(utils.AesCbc256HmacSha256B64, signPrivateKeyPem, symmetricKeyAesKey, symmetricKeyHashKey)
	if err != nil {
		logger.Errorf("[%v] encryption signing private key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption signing private key failed")
	}
	protectedTLSPrivateKey, err := utils.Encrypt(utils.AesCbc256HmacSha256B64, tlsPrivateKeyPem, symmetricKeyAesKey, symmetricKeyHashKey)
	if err != nil {
		logger.Errorf("[%v] encryption tls private key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption tls private key failed")
	}
	protectedRSAPrivateKey, err := utils.Encrypt(utils.AesCbc256HmacSha256B64, rsaPrivateKeyPem, symmetricKeyAesKey, symmetricKeyHashKey)
	if err != nil {
		logger.Errorf("[%v] encryption rsa private key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption rsa private key failed")
	}
	u.ProtectedSignPrivateKey = protectedSignPrivateKey
	u.ProtectedTLSPrivateKey = protectedTLSPrivateKey
	u.ProtectedRSAPrivateKey = protectedRSAPrivateKey

	pubKey, _ := signPrivateKey.PublicKey()
	pubKeyPem, err := pubKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the signing public key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the signing public key to pem format")
	}
	u.SignPublicKey = string(pubKeyPem)

	pubKey, _ = tlsPrivateKey.PublicKey()
	pubKeyPem, err = pubKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the tls public key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the tls public key to pem format")
	}
	u.TLSPublicKey = string(pubKeyPem)

	pubKey, _ = rsaPrivateKey.PublicKey()
	pubKeyPem, err = pubKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the rsa public key to pem format error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the rsa public key to pem format")
	}
	u.RSAPublicKey = string(pubKeyPem)

	stretchedKeyAesKey, err := factory.CryptoKeyImport(stretchedKey.Enc, crypto.AesCbc256)
	if err != nil {
		logger.Errorf("[%v] import aes password error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import aes password key")
	}
	stretchedKeyHashKey, err := factory.CryptoKeyImport(stretchedKey.Mac, crypto.HmacSha256)
	if err != nil {
		logger.Errorf("[%v] import hash password error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import hash password key")
	}
	protectedSymmetricKey, err := utils.Encrypt(utils.AesCbc256HmacSha256B64, symmetricKey.Key(), stretchedKeyAesKey, stretchedKeyHashKey)
	if err != nil {
		logger.Errorf("[%v] encryption symmetric key error: %v", u.UserID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption symmetric key failed")
	}
	u.ProtectedSymmetricKey = protectedSymmetricKey

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
