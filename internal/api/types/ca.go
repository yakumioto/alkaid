/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package types

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"

	"github.com/pkg/errors"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

const (
	SignCAType = "sign"
	TLSCAType  = "tls"
)

var (
	typeCAMaps = map[string]string{
		SignCAType: SignCAType,
		TLSCAType:  TLSCAType,
	}
)

func chekcType(typ string) bool {
	_, ok := typeCAMaps[typ]
	return ok
}

// CA Sign CA or TLS CA
type CA struct {
	ID             int64  `json:"-"`
	OrganizationID string `json:"organization_id,omitempty"`
	Type           string `json:"type,omitempty"`
	PrivateKey     []byte `json:"-"`
	Certificate    []byte `json:"certificate,omitempty"`
	CreateAt       int64  `json:"create_at,omitempty"`
	UpdateAt       int64  `json:"update_at,omitempty"`
}

// NewCA New CA
func NewCA(org *Organization, typ string) (*CA, error) {
	if !chekcType(typ) {
		return nil, errors.New("error type")
	}

	if org == nil || org.OrganizationID == "" {
		return nil, errors.New("error organization")
	}

	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return nil, err
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

	// // set the organization for the subject
	template.Subject = crypto.SubjectTemplateAdditional(
		org.OrganizationID,
		fmt.Sprintf("%s.%s", "ca", org.Domain),
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
		return nil, err
	}

	privPemBytes, err := crypto.PrivateKeyExport(priv)
	if err != nil {
		return nil, err
	}

	x509PemBytes := crypto.X509Export(x509Cert)

	return &CA{
		OrganizationID: org.OrganizationID,
		Type:           typ,
		PrivateKey:     privPemBytes,
		Certificate:    x509PemBytes,
	}, nil
}

// SignCertificate creates a signed certificate based on a built-in template
func (c *CA) SignCertificate(
	org *Organization,
	orgUnits,
	alternateNames []string,
	pub *ecdsa.PublicKey,
) (*x509.Certificate, error) {
	if org == nil || org.OrganizationID == "" {
		return nil, errors.New("error organization")
	}

	template := crypto.X509Template()
	switch c.Type {
	case TLSCAType:
		template.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
		template.ExtKeyUsage = []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		}
	case SignCAType:
		template.KeyUsage = x509.KeyUsageDigitalSignature
	}

	subject := crypto.SubjectTemplateAdditional(
		org.OrganizationID,
		org.Name,
		org.Country,
		org.Province,
		org.Locality,
		org.OrganizationalUnit,
		org.StreetAddress,
		org.PostalCode,
	)
	subject.OrganizationalUnit = append(subject.OrganizationalUnit, orgUnits...)

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

	cert, err := c.SignCert()
	if err != nil {
		return nil, err
	}

	priv, err := c.Signer()
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
func (c *CA) SignCert() (*x509.Certificate, error) {
	var cert *x509.Certificate
	var err error

	block, _ := pem.Decode(c.Certificate)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.Errorf("bytes are not PEM encoded")
	}
	cert, err = x509.ParseCertificate(block.Bytes)

	return cert, err
}

func (c *CA) Signer() (*crypto.ECDSASigner, error) {
	block, _ := pem.Decode(c.PrivateKey)
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
