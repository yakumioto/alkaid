/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package certificate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/third_party/github.com/hyperledger/fabric/common/crypto"
)

func TestNewCA(t *testing.T) {
	_, _, err := NewCA(types.NewOrganization(), "1")

	assert.NoError(t, err, "new ca error")
}

func TestSignCertificate(t *testing.T) {
	priv, cert, _ := NewCA(types.NewOrganization(), "1")
	privByte, _ := crypto.PrivateKeyExport(priv)
	certByte := crypto.X509Export(cert)

	pk, _ := crypto.GeneratePrivateKey()
	_, err := SignCertificate(types.NewOrganization(), "1", "orderer", nil, &pk.PublicKey, privByte, certByte)

	assert.NoError(t, err, "sign certificate error")
}
