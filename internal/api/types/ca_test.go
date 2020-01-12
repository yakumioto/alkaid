/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func TestNewCA(t *testing.T) {
	// error type test
	_, err := NewCA(&Organization{OrganizationID: "testorg"}, "type")
	assert.EqualError(t, err, "error type")

	// empty organization id test
	_, err = NewCA(&Organization{}, SignCAType)
	assert.Error(t, err, "error organization")

	_, err = NewCA(&Organization{OrganizationID: "testorg"}, TLSCAType)
	assert.NoError(t, err)
}

func TestCA_SignCertificate(t *testing.T) {
	ca, _ := NewCA(&Organization{OrganizationID: "testorg"}, SignCAType)

	// empty organization test
	_, err := ca.SignCertificate(nil, []string{}, nil, nil)
	assert.EqualError(t, err, "error organization")

	priv, _ := crypto.GeneratePrivateKey()
	_, err = ca.SignCertificate(
		&Organization{OrganizationID: "testorg", Domain: "testorg.com"},
		[]string{OrdererMSPType},
		[]string{"mioto.me"},
		&priv.PublicKey)
	assert.NoError(t, err)
}

func TestCA_SignCert(t *testing.T) {
	ca, _ := NewCA(&Organization{OrganizationID: "testorg"}, SignCAType)
	_, err := ca.SignCert()
	assert.NoError(t, err)
}

func TestCA_Signer(t *testing.T) {
	ca, _ := NewCA(&Organization{OrganizationID: "testorg"}, SignCAType)
	_, err := ca.Signer()
	assert.NoError(t, err)
}
