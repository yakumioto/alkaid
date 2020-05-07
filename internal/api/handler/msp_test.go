/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package handler

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/yakumioto/alkaid/internal/api/errors"
	"github.com/yakumioto/alkaid/internal/api/types"
)

var (
	testCreateMSP = organization + organizationID + "/msp"
	testGetMSP    = organization + organizationID + "/msp/:userID"
)

func setMSPAllFields(msp *types.MSP) {
	msp.UserID = "user1_id"
	msp.Name = "user1_name"
	msp.Type = "admin"
}

func TestCreateMSP(t *testing.T) {
	r := testInit()

	r.POST(organization, CreateOrganization)
	r.POST(testCreateMSP, CreateMSP)

	expectedBody := fmt.Sprintln(errors.NewErrors(errors.BadRequestData).Error())

	msp := types.NewMSP()

	// empty test
	setMSPAllFields(msp)
	msp.UserID = ""
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusBadRequest, expectedBody)

	setMSPAllFields(msp)
	msp.Name = ""
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusBadRequest, expectedBody)

	setMSPAllFields(msp)
	msp.Type = ""
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusBadRequest, expectedBody)

	// orgnaization not exist test
	expectedBody = fmt.Sprintln(errors.NewErrors(errors.DataNotExists).Error())

	setMSPAllFields(msp)
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusNotFound, expectedBody)

	// ca not exist test
	org := types.NewOrganization()
	setOrganizationAllFields(org)
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusOK, nil)
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusOK, nil)

	// exist test
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusBadRequest, nil)
}

func TestGetMSPByUserID(t *testing.T) {
	r := testInit()

	r.POST(organization, CreateOrganization)
	r.POST(testCreateMSP, CreateMSP)
	r.GET(testGetMSP, GetMSPByUserID)

	// no organization test
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id/msp/user1_id", nil, http.StatusNotFound, nil)

	// no msp test
	org := types.NewOrganization()
	setOrganizationAllFields(org)
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusOK, nil)
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id/msp/user1_id", nil, http.StatusNotFound, nil)

	// ok test
	msp := types.NewMSP()
	setMSPAllFields(msp)
	testHTTPEqual(t, r, http.MethodPost, organization+"/org1_id/msp", msp, http.StatusOK, nil)
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id/msp/user1_id", nil, http.StatusOK, nil)
}
