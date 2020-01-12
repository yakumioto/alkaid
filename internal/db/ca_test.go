/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func insertTestCA(t *testing.T) *CA {
	org := insertTestOrganization(t)

	ca := &CA{OrganizationID: org.OrganizationID, Type: "orderer"}
	err := CreateCA(ca)
	assert.NoError(t, err)

	return ca
}

func TestCreateCA(t *testing.T) {
	testInit()

	// empty organization test
	err := CreateCA(&CA{OrganizationID: "testorg"})
	assert.EqualError(t, err, "organization not exists [organization_id: testorg]")

	// no err test
	ca := insertTestCA(t)

	// exist test
	err = CreateCA(ca)
	assert.EqualError(t, err, "ca already exists [organization_id: testorg, type: orderer]")
}

func TestQueryCAByOrganizationIDAndType(t *testing.T) {
	testInit()

	_, err := QueryCAByOrganizationIDAndType("testnotexist", "orderer")
	assert.EqualError(t, err, "ca not exists [organization_id: testnotexist, type: orderer]")

	ca := insertTestCA(t)

	_, err = QueryCAByOrganizationIDAndType(ca.OrganizationID, ca.Type)
	assert.NoError(t, err)
}
