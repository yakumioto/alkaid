/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// compute Subject Key Identifier
func ComputeSKI(privKey *ecdsa.PrivateKey) []byte {
	// Marshall the public key
	raw := elliptic.Marshal(privKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)

	// Hash it
	hash := sha256.Sum256(raw)
	return hash[:]
}

// SubjectTemplateAdditional Additional for X509 subject
func SubjectTemplateAdditional(
	orgName,
	commonName,
	country,
	province,
	locality,
	orgUnit,
	streetAddress,
	postalCode string,
) pkix.Name {
	var name pkix.Name
	if len(orgName) >= 1 {
		name.Organization = []string{orgName}
	}
	if len(commonName) >= 1 {
		name.CommonName = commonName
	}

	if len(country) >= 1 {
		name.Country = []string{country}
	}
	if len(province) >= 1 {
		name.Province = []string{province}
	}

	if len(locality) >= 1 {
		name.Locality = []string{locality}
	}
	if len(orgUnit) >= 1 {
		name.OrganizationalUnit = []string{orgUnit}
	}
	if len(streetAddress) >= 1 {
		name.StreetAddress = []string{streetAddress}
	}
	if len(postalCode) >= 1 {
		name.PostalCode = []string{postalCode}
	}
	return name
}

// X509Template default template for X509 certificates
func X509Template() x509.Certificate {
	// generate a serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	// set expiry to around 10 years
	expiry := 3650 * 24 * time.Hour
	// round minute and backdate 5 minutes
	notBefore := time.Now().Round(time.Minute).Add(-5 * time.Minute).UTC()

	// basic template to use
	X509 := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notBefore.Add(expiry).UTC(),
		BasicConstraintsValid: true,
	}
	return X509
}

// GenCertificateECDSA generate a signed X509 certificate using ECDSA
func GenCertificateECDSA(
	template,
	parent *x509.Certificate,
	pub *ecdsa.PublicKey,
	priv interface{},
) (*x509.Certificate, error) {
	// create the x509 public cert
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}
	return x509Cert, nil
}
