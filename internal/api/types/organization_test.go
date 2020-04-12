/*
 *  Copyright 2020. The Alkaid Authors. All rights reserved.
 *  Use of this source code is governed by a MIT-style
 *  license that can be found in the LICENSE file.
 *  Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package types

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yakumioto/alkaid/internal/config"
)

func testNewOrganization() *Organization {
	org := NewOrganization()
	org.OrganizationID = "testorg"

	return org
}

func TestOrganization_CreateMSPDir(t *testing.T) {
	org := testNewOrganization()
	signca, _ := NewCA(org, SignCAType)
	tlsca, _ := NewCA(org, TLSCAType)

	config.FileSystemPath = path.Join(os.TempDir(), "alkaid")
	err := org.CreateMSPDir(signca, tlsca)
	assert.NoError(t, err)

	_, err = os.Stat(path.Join(config.FileSystemPath, "testorg", "msp", "cacerts", "cert.pem"))
	assert.NoError(t, err)

	_, err = os.Stat(path.Join(config.FileSystemPath, "testorg", "msp", "tlscacerts", "cert.pem"))
	assert.NoError(t, err)

	_ = os.RemoveAll(path.Join(config.FileSystemPath))
}
