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

func setNetworkAllFields(network *types.Network) {
	network.NetworkID = "test_network"
	network.Type = "docker"
}

func TestCreateNetwork(t *testing.T) {
	r := testInit()

	r.POST("/network", CreateNetwork)

	network := types.NewNetwork()
	// empty test
	expectedBody := fmt.Sprintln(errors.NewErrors(errors.BadAuthenticationData).Error())

	setNetworkAllFields(network)
	network.NetworkID = ""
	testHTTPEqual(t, r, http.MethodPost, "/network", network, http.StatusBadRequest, expectedBody)

	setNetworkAllFields(network)
	network.Type = ""
	testHTTPEqual(t, r, http.MethodPost, "/network", network, http.StatusBadRequest, expectedBody)

	// ok test
	setNetworkAllFields(network)
	testHTTPEqual(t, r, http.MethodPost, "/network", network, http.StatusOK, nil)

	// exist test
	expectedBody = fmt.Sprintln(errors.NewErrors(errors.DataAlreadyExists).Error())
	testHTTPEqual(t, r, http.MethodPost, "/network", network, http.StatusBadRequest, expectedBody)
}

func TestGetNetworkByID(t *testing.T) {
	r := testInit()

	r.POST("/network", CreateNetwork)
	r.GET("/network/:networkID", GetNetworkByID)

	network := types.NewNetwork()
	setNetworkAllFields(network)
	testHTTPEqual(t, r, http.MethodPost, "/network", network, http.StatusOK, nil)

	// not exist test
	expectedBody := fmt.Sprintln(errors.NewErrors(errors.DataNotExists).Error())
	testHTTPEqual(t, r, http.MethodGet, "/network/not_exist", nil, http.StatusBadRequest, expectedBody)

	// ok test
	testHTTPEqual(t, r, http.MethodGet, "/network/test_network", nil, http.StatusOK, nil)
}
