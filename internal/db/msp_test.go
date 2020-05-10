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

func insertTestMSP(t *testing.T) *User {
	org := insertTestOrganization(t)

	msp := &User{OrganizationID: org.OrganizationID, UserID: "testmsp"}
	err := CreateMSP(msp)
	assert.NoError(t, err)

	return msp
}

func TestCreateMSP(t *testing.T) {
	testInit()

	// empty organization test
	err := CreateMSP(&User{OrganizationID: "testorg", UserID: "testmsp"})
	assert.EqualError(t, err, "organization not exists [organization_id: testorg]")

	// no err test
	msp := insertTestMSP(t)

	// exist test
	err = CreateMSP(msp)
	assert.EqualError(t, err, "user already exists [organization_id: testorg, user_id: testmsp]")
}

func TestQueryMSPByOrganizationIDAndMSPID(t *testing.T) {
	testInit()

	_, err := QueryMSPByOrganizationIDAndUserID("testnotexist", "testmsp")
	assert.EqualError(t, err, "user not exists [organization_id: testnotexist, user_id: testmsp]")

	msp := insertTestMSP(t)

	_, err = QueryMSPByOrganizationIDAndUserID(msp.OrganizationID, msp.UserID)
	assert.NoError(t, err)
}
