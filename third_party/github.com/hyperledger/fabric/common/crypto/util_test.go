/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyExport(t *testing.T) {
	pemPriv := []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgtyYl8B39JQKyxntg
rOSxxEgPRXOJ8eomRQ9Lbtf+nYOhRANCAAQeFAMwcw0maiH2ucEZYGnJKCijdK/7
Kvcz6Cot90aldVFxrkSxrEO4cvJglqkKenTJAkGGM3RFj8GMxdewQCT/
-----END PRIVATE KEY-----
`)
	block, _ := pem.Decode(pemPriv)

	key, _ := x509.ParsePKCS8PrivateKey(block.Bytes)

	priv, _ := key.(*ecdsa.PrivateKey)

	privByte, err := PrivateKeyExport(priv)
	assert.NoError(t, err)
	assert.Equal(t, pemPriv, privByte)
}

func TestX509Export(t *testing.T) {
	pemX509 := []byte(`-----BEGIN CERTIFICATE-----
MIICUTCCAfegAwIBAgIQZSSg6qZ22Gzaj2bYj4frVzAKBggqhkjOPQQDAjBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMDAyMDMxMzI3MDBaFw0zMDAxMzExMzI3MDBa
MHMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMRwwGgYDVQQD
ExNjYS5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
HhQDMHMNJmoh9rnBGWBpySgoo3Sv+yr3M+gqLfdGpXVRca5EsaxDuHLyYJapCnp0
yQJBhjN0RY/BjMXXsEAk/6NtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1UdJQQWMBQG
CCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdDgQiBCBq
iukiCVFIGV4lyp+HnNhk6hnEKybH71xLTTn2bCtvRTAKBggqhkjOPQQDAgNIADBF
AiEAn7zAHnpHdSDY8gI6GjdYhM0JE2mcbJt3aDnyBjF4bk8CIHu548lYeI+UBDEH
nBwbCVBediP7pG1GkMWnZC54RgOH
-----END CERTIFICATE-----
`)

	block, _ := pem.Decode(pemX509)
	cert, _ := x509.ParseCertificate(block.Bytes)

	assert.Equal(t, pemX509, X509Export(cert))
}
