/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func PrivateKeyExport(priv *ecdsa.PrivateKey) ([]byte, error) {
	pkcs8Encoded, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal private key")
	}

	return pemExport("PRIVATE KEY", pkcs8Encoded), nil
}

func X509Export(cert *x509.Certificate) []byte {
	return pemExport("CERTIFICATE", cert.Raw)
}

func pemExport(pemType string, raw []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: pemType, Bytes: raw})
}
