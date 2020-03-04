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

	"github.com/yakumioto/alkaid/internal/api/types"
)

func insertTestNetwork(t *testing.T) *Network {
	network := &Network{
		NetworkID:       "testnetwork",
		Type:            types.DockerNetworkType,
		DockerNetworkID: "testdockernetworkid",
	}
	err := CreateNetwork(network)
	assert.NoError(t, err)

	return network
}

func TestCreateNetwork(t *testing.T) {
	testInit()

	err := CreateNetwork(&Network{Type: types.DockerNetworkType})
	assert.NoError(t, err)

	// no err test
	network := insertTestNetwork(t)

	// exist test
	err = CreateNetwork(network)
	assert.EqualError(t, err, "network already exists [network_id: testnetwork]")
}

func TestQueryNetworkByNetworkID(t *testing.T) {
	testInit()

	_, err := QueryNetworkByNetworkID("testnotexist")
	assert.EqualError(t, err, "network not exists [network_id: testnotexist]")

	network := insertTestNetwork(t)

	_, err = QueryNetworkByNetworkID(network.NetworkID)
	assert.NoError(t, err)
}
