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
	"net/http"
	"testing"

	"github.com/yakumioto/alkaid/internal/api/types"
)

var (
	signCA = organization + organizationID + "/signca"
	tlsCA  = organization + organizationID + "/tlsca"
)

func TestGetCAByOrganizationID(t *testing.T) {
	r := testInit()

	r.POST(organization, CreateOrganization)
	r.GET(signCA, GetCAByOrganizationID)
	r.GET(tlsCA, GetCAByOrganizationID)

	// no organization test
	testHTTPEqual(t, r, http.MethodGet, organization+"/test/signca", nil, http.StatusBadRequest, nil)

	org := types.NewOrganization()
	setOrganizationAllFields(org)
	testHTTPEqual(t, r, http.MethodPost, organization, org, http.StatusOK, nil)

	// ok test
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id/signca", nil, http.StatusOK, nil)
	testHTTPEqual(t, r, http.MethodGet, organization+"/org1_id/tlsca", nil, http.StatusOK, nil)
}
