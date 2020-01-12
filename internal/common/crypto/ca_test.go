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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeSKI(t *testing.T) {
	priv, _ := GeneratePrivateKey()
	ski := ComputeSKI(priv)
	assert.NotNil(t, ski)
}

func TestSubjectTemplateAdditional(t *testing.T) {
	subject := SubjectTemplateAdditional(
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
	)
	assert.NotNil(t, subject)
}

func TestX509Template(t *testing.T) {
	template := X509Template()
	assert.NotNil(t, template)
}

func TestGenCertificateECDSA(t *testing.T) {
	priv, _ := GeneratePrivateKey()
	signer := &ECDSASigner{PrivateKey: priv}
	template := X509Template()
	template.Subject = SubjectTemplateAdditional(
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
	)
	_, err := GenCertificateECDSA(&template, &template, &priv.PublicKey, signer)
	assert.NoError(t, err)
}
