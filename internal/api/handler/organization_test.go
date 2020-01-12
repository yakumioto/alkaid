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
	organization   = "/organization"
	organizationID = "/:organizationID"
)

func setOrganizationAllFields(org *types.Organization) {
	org.OrganizationID = "org1_id"
	org.Name = "org1_name"
	org.Domain = "org1.com"
	org.Type = "orderer"
}

func TestCreateOrganization(t *testing.T) {
	r := testInit()

	org := types.NewOrganization()

	r.POST(organization, CreateOrganization)

	expectedBody := fmt.Sprintln(errors.NewErrors(errors.BadAuthenticationData).Error())

	// empty test
	setOrganizationAllFields(org)
	org.OrganizationID = ""
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	setOrganizationAllFields(org)
	org.Name = ""
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	setOrganizationAllFields(org)
	org.Domain = ""
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	setOrganizationAllFields(org)
	org.Type = ""
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	// fqdn test
	setOrganizationAllFields(org)
	org.Domain = "domain"
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	// oneof test
	setOrganizationAllFields(org)
	org.Type = "type"
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)

	// ok test
	setOrganizationAllFields(org)
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusOK, nil)

	// exist test
	expectedBody = fmt.Sprintln(errors.NewErrors(errors.DataAlreadyExists).Error())
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusBadRequest, expectedBody)
}

func TestGetOrganization(t *testing.T) {
	r := testInit()

	r.POST(organization, CreateOrganization)
	r.GET(organization+organizationID, GetOrganizationByID)

	org := types.NewOrganization()
	setOrganizationAllFields(org)
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusOK, nil)

	// not exist test
	expectedBody := fmt.Sprintln(errors.NewErrors(errors.DataNotExists).Error())
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_not_exists", nil, http.StatusBadRequest, expectedBody)

	// ok test
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id", nil, http.StatusOK, nil)
}
