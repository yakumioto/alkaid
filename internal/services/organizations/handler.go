/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package organizations

import (
	"encoding/base64"
	"net/http"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/crypto/factory"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/errors"
)

var (
	logger = log.GetPackageLogger("services.organizations")
)

type CreateRequest struct {
	OrganizationID      string `json:"organizationId,omitempty" validate:"required"`
	Name                string `json:"name,omitempty" validate:"required"`
	Domain              string `json:"domain,omitempty" validate:"required,fqdn"`
	Description         string `json:"description,omitempty"`
	Country             string `json:"country,omitempty"`
	Province            string `json:"province,omitempty"`
	Locality            string `json:"locality,omitempty"`
	OrganizationalUnit  string `json:"organizationalUnit,omitempty"`
	StreetAddress       string `json:"streetAddress,omitempty"`
	PostalCode          string `json:"postalCode,omitempty"`
	TransactionPassword string `json:"transactionPassword" validate:"required"` // 交易密码仅用来加解密 PrivateKey
}

func Create(req *CreateRequest) (*Organization, error) {
	org := newOrganizationByCreateRequest(req)

	signCAPrivateKey, err := factory.CryptoKeyGen(crypto.EcdsaP256)
	if err != nil {
		logger.Errorf("[%v] generate signature key error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate signature key")
	}
	signCAPrivateKeyPem, err := signCAPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the signature key to pem format error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the signature key to pem format")
	}

	tlsCAPrivateKey, err := factory.CryptoKeyGen(crypto.EcdsaP256)
	if err != nil {
		logger.Errorf("[%v] generate tls key error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to generate tls key")
	}
	tlsCAPrivateKeyPem, err := tlsCAPrivateKey.Bytes()
	if err != nil {
		logger.Errorf("[%v] convert the tls key to pem format error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to convert the tls key to pem format")
	}

	aesKey, err := factory.CryptoKeyImport([]byte(req.TransactionPassword), crypto.AesCbc256)
	if err != nil {
		logger.Errorf("[%v] import transaction password error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"failed to import transaction password key")
	}
	protectedSignPrivateKey, err := aesKey.Encrypt(signCAPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption signing key error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption signing key failed")
	}
	protectedTLSPrivateKey, err := aesKey.Encrypt(tlsCAPrivateKeyPem)
	if err != nil {
		logger.Errorf("[%v] encryption tls key error: %v", req.OrganizationID, err)
		return nil, errors.NewError(http.StatusInternalServerError, errors.ErrServerUnknownError,
			"encryption tls key failed")
	}

	org.ProtectedSignCAPrivateKey = base64.StdEncoding.EncodeToString(protectedSignPrivateKey)
	org.ProtectedTLSCAPrivateKey = base64.StdEncoding.EncodeToString(protectedTLSPrivateKey)

	// TODO: 生成证书

	return nil, nil
}
