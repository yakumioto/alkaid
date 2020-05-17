/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package certificate

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"net"

	"github.com/pkg/errors"

	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

type PkixName struct {
	OrgName       string
	CommonName    string
	Country       string
	Province      string
	Locality      string
	OrgUnit       string
	StreetAddress string
	PostalCode    string
}

func NewCA(org *types.Organization, commonName string) (*ecdsa.PrivateKey, *x509.Certificate, error) {
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil, nil, err
	}

	template := crypto.X509Template()

	// this is a CA
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageDigitalSignature |
		x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
		x509.KeyUsageCRLSign
	template.ExtKeyUsage = []x509.ExtKeyUsage{
		x509.ExtKeyUsageClientAuth,
		x509.ExtKeyUsageServerAuth,
	}

	// set the organization for the subject
	template.Subject = crypto.SubjectTemplateAdditional(
		org.Domain,
		commonName,
		org.Country,
		org.Province,
		org.Locality,
		org.OrganizationalUnit,
		org.StreetAddress,
		org.PostalCode,
	)
	template.SubjectKeyId = crypto.ComputeSKI(priv)

	x509Cert, err := crypto.GenCertificateECDSA(
		&template,
		&template,
		&priv.PublicKey,
		priv,
	)
	if err != nil {
		return nil, nil, err
	}

	return priv, x509Cert, nil
}

func SignCertificate(
	org *types.Organization,
	commonName,
	orgUnits string,
	alternateNames []string,
	pub *ecdsa.PublicKey,
	caPrivByte,
	caCertificateByte []byte) (*x509.Certificate, error) {
	template := crypto.X509Template()
	switch orgUnits {
	case types.OrdererMSPType, types.PeerMSPType:
		template.KeyUsage = x509.KeyUsageDigitalSignature
	case types.AdminMSPType, types.ClientMSPType:
		template.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		template.ExtKeyUsage = []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}
	}

	subject := crypto.SubjectTemplateAdditional(
		"",
		commonName,
		org.Country,
		org.Province,
		org.Locality,
		org.OrganizationalUnit,
		org.StreetAddress,
		org.PostalCode,
	)
	subject.OrganizationalUnit = append(subject.OrganizationalUnit, orgUnits)

	template.Subject = subject

	for _, san := range alternateNames {
		// try to parse as an IP address first
		ip := net.ParseIP(san)
		if ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, san)
		}
	}

	var (
		priv *crypto.ECDSASigner
		cert *x509.Certificate
		err  error
	)

	priv, err = Signer(caPrivByte)
	if err != nil {
		return nil, err
	}
	cert, err = SignCert(caCertificateByte)
	if err != nil {
		return nil, err
	}

	return crypto.GenCertificateECDSA(
		&template,
		cert,
		pub,
		priv)
}

// SignCert load a ecdsa cert from Certificate
func SignCert(certByte []byte) (*x509.Certificate, error) {
	var cert *x509.Certificate
	var err error

	block, _ := pem.Decode(certByte)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.Errorf("bytes are not PEM encoded")
	}
	cert, err = x509.ParseCertificate(block.Bytes)

	return cert, err
}

func Signer(privKey []byte) (*crypto.ECDSASigner, error) {
	block, _ := pem.Decode(privKey)
	if block == nil {
		return nil, errors.New("bytes are not PEM encoded")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.WithMessage(err, "pem bytes are not PKCS8 encoded ")
	}

	priv, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("pem bytes do not contain an EC private key")
	}
	return &crypto.ECDSASigner{PrivateKey: priv}, nil
}
