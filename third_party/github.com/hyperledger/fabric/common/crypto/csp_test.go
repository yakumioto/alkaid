/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	priv, err := GeneratePrivateKey()
	assert.NoError(t, err, "Failed to generate private key")
	assert.NotNil(t, priv, "Should have returned an *ecdsa.Key")
	assert.NoError(t, err,
		"Expected to find private key file")
}

func TestECDSASigner(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate private key: %s", err)
	}

	signer := ECDSASigner{
		PrivateKey: priv,
	}
	assert.Equal(t, priv.Public(), signer.Public().(*ecdsa.PublicKey))
	digest := []byte{1}
	sig, err := signer.Sign(rand.Reader, digest, nil)
	if err != nil {
		t.Fatalf("Failed to create signature: %s", err)
	}

	// unmarshal signature
	ecdsaSig := &ECDSASignature{}
	_, err = asn1.Unmarshal(sig, ecdsaSig)
	if err != nil {
		t.Fatalf("Failed to unmarshal signature: %s", err)
	}
	// s should not be greater than half order of curve
	halfOrder := new(big.Int).Div(priv.PublicKey.Curve.Params().N, big.NewInt(2))

	if ecdsaSig.S.Cmp(halfOrder) == 1 {
		t.Error("Expected signature with Low S")
	}

	// ensure signature is valid by using standard verify function
	ok := ecdsa.Verify(&priv.PublicKey, digest, ecdsaSig.R, ecdsaSig.S)
	assert.True(t, ok, "Expected valid signature")
}
