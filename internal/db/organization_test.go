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

func insertTestOrganization(t *testing.T) *Organization {
	org := &Organization{OrganizationID: "testorg", Domain: "testorg.com"}
	err := CreateOrganization(org)
	assert.NoError(t, err)

	return org
}

func TestCreateOrganization(t *testing.T) {
	testInit()

	// empty test
	err := CreateOrganization(&Organization{})
	assert.NoError(t, err)

	// no err test
	org := insertTestOrganization(t)

	// exist test
	err = CreateOrganization(org)
	assert.EqualError(t, err, "organization already exists [organization_id: testorg]")
}

func TestQueryOrganizationByOrgID(t *testing.T) {
	testInit()

	// not exist test
	_, err := QueryOrganizationByOrgID("testnotexist")
	assert.EqualError(t, err, "organization not exists [organization_id: testnotexist]")

	org := insertTestOrganization(t)

	_, err = QueryOrganizationByOrgID(org.OrganizationID)
	assert.NoError(t, err)
}
