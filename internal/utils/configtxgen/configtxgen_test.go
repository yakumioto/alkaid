/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package configtxgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGetOrganization() *OrganizationConf {
	return &OrganizationConf{
		Name:          "TestOrg",
		MSPIdentifier: "TestMSP",
		AdminCert: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNYRENDQWdPZ0F3SUJBZ0lSQUl0cmk5U051aEc0TEVzN1BYM" +
			"G1GVzR3Q2dZSUtvWkl6ajBFQXdJd2JERUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQ" +
			"mdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhGREFTQmdOVkJBb1RDMlY0WVcxd2JHVXVZMjl0TVJvd0dBWURWUVFERXhGMGJIT" +
			"mpZUzVsCmVHRnRjR3hsTG1OdmJUQWVGdzB5TURBME1EY3dPRE00TURCYUZ3MHpNREEwTURVd09ETTRNREJhTUZreEN6QUoKQmdOV" +
			"kJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVRZVzRnUm5KaApibU5wYzJOdk1SMHdHd" +
			"1lEVlFRREV4UnZjbVJsY21WeU5DNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkx1N" +
			"G4wK2plaUFwVVNrWkh4R0wvMFg1UUdCUitmNU1MWHhoSE5OUlR3anUKcDdGazlFcmZwQjd2ZlVLdjNzM1dsN2tWeUpTZ0hPMWNVc" +
			"mxlM3RkTi9UeWpnWmd3Z1pVd0RnWURWUjBQQVFILwpCQVFEQWdXZ01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUJCZ2dyQmdFR" +
			"kJRY0RBakFNQmdOVkhSTUJBZjhFCkFqQUFNQ3NHQTFVZEl3UWtNQ0tBSUNvNGJVcmxGekZtR085TmREemhKeWQ3T1YxdFVvb20wR" +
			"Ss3d2k5R3JZRkcKTUNrR0ExVWRFUVFpTUNDQ0ZHOXlaR1Z5WlhJMExtVjRZVzF3YkdVdVkyOXRnZ2h2Y21SbGNtVnlOREFLQmdnc" +
			"Qpoa2pPUFFRREFnTkhBREJFQWlCQmZBbDA2TmJ4VVhCdDQ1M01UdEhOa1NvZ05YZHVQa1dUdC9TTUVkY2Q1Z0lnClI0ZFExNVFTR" +
			"2o0ZENUWVRyT1ZkUFcwa3ZyS3dueExCTzVoUTk3Q1hGbk09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
		SignRootCert: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNYRENDQWdPZ0F3SUJBZ0lSQUl0cmk5U051aEc0TEVzN1BYM" +
			"G1GVzR3Q2dZSUtvWkl6ajBFQXdJd2JERUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQ" +
			"mdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhGREFTQmdOVkJBb1RDMlY0WVcxd2JHVXVZMjl0TVJvd0dBWURWUVFERXhGMGJIT" +
			"mpZUzVsCmVHRnRjR3hsTG1OdmJUQWVGdzB5TURBME1EY3dPRE00TURCYUZ3MHpNREEwTURVd09ETTRNREJhTUZreEN6QUoKQmdOV" +
			"kJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVRZVzRnUm5KaApibU5wYzJOdk1SMHdHd" +
			"1lEVlFRREV4UnZjbVJsY21WeU5DNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkx1N" +
			"G4wK2plaUFwVVNrWkh4R0wvMFg1UUdCUitmNU1MWHhoSE5OUlR3anUKcDdGazlFcmZwQjd2ZlVLdjNzM1dsN2tWeUpTZ0hPMWNVc" +
			"mxlM3RkTi9UeWpnWmd3Z1pVd0RnWURWUjBQQVFILwpCQVFEQWdXZ01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUJCZ2dyQmdFR" +
			"kJRY0RBakFNQmdOVkhSTUJBZjhFCkFqQUFNQ3NHQTFVZEl3UWtNQ0tBSUNvNGJVcmxGekZtR085TmREemhKeWQ3T1YxdFVvb20wR" +
			"Ss3d2k5R3JZRkcKTUNrR0ExVWRFUVFpTUNDQ0ZHOXlaR1Z5WlhJMExtVjRZVzF3YkdVdVkyOXRnZ2h2Y21SbGNtVnlOREFLQmdnc" +
			"Qpoa2pPUFFRREFnTkhBREJFQWlCQmZBbDA2TmJ4VVhCdDQ1M01UdEhOa1NvZ05YZHVQa1dUdC9TTUVkY2Q1Z0lnClI0ZFExNVFTR" +
			"2o0ZENUWVRyT1ZkUFcwa3ZyS3dueExCTzVoUTk3Q1hGbk09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
		TLSRootCert: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNYRENDQWdPZ0F3SUJBZ0lSQUl0cmk5U051aEc0TEVzN1BYM" +
			"G1GVzR3Q2dZSUtvWkl6ajBFQXdJd2JERUwKTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmhiR2xtYjNKdWFXRXhGakFVQ" +
			"mdOVkJBY1REVk5oYmlCRwpjbUZ1WTJselkyOHhGREFTQmdOVkJBb1RDMlY0WVcxd2JHVXVZMjl0TVJvd0dBWURWUVFERXhGMGJIT" +
			"mpZUzVsCmVHRnRjR3hsTG1OdmJUQWVGdzB5TURBME1EY3dPRE00TURCYUZ3MHpNREEwTURVd09ETTRNREJhTUZreEN6QUoKQmdOV" +
			"kJBWVRBbFZUTVJNd0VRWURWUVFJRXdwRFlXeHBabTl5Ym1saE1SWXdGQVlEVlFRSEV3MVRZVzRnUm5KaApibU5wYzJOdk1SMHdHd" +
			"1lEVlFRREV4UnZjbVJsY21WeU5DNWxlR0Z0Y0d4bExtTnZiVEJaTUJNR0J5cUdTTTQ5CkFnRUdDQ3FHU000OUF3RUhBMElBQkx1N" +
			"G4wK2plaUFwVVNrWkh4R0wvMFg1UUdCUitmNU1MWHhoSE5OUlR3anUKcDdGazlFcmZwQjd2ZlVLdjNzM1dsN2tWeUpTZ0hPMWNVc" +
			"mxlM3RkTi9UeWpnWmd3Z1pVd0RnWURWUjBQQVFILwpCQVFEQWdXZ01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUJCZ2dyQmdFR" +
			"kJRY0RBakFNQmdOVkhSTUJBZjhFCkFqQUFNQ3NHQTFVZEl3UWtNQ0tBSUNvNGJVcmxGekZtR085TmREemhKeWQ3T1YxdFVvb20wR" +
			"Ss3d2k5R3JZRkcKTUNrR0ExVWRFUVFpTUNDQ0ZHOXlaR1Z5WlhJMExtVjRZVzF3YkdVdVkyOXRnZ2h2Y21SbGNtVnlOREFLQmdnc" +
			"Qpoa2pPUFFRREFnTkhBREJFQWlCQmZBbDA2TmJ4VVhCdDQ1M01UdEhOa1NvZ05YZHVQa1dUdC9TTUVkY2Q1Z0lnClI0ZFExNVFTR" +
			"2o0ZENUWVRyT1ZkUFcwa3ZyS3dueExCTzVoUTk3Q1hGbk09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
	}
}

func TestGetGenesisBlock(t *testing.T) {
	conf := &GenesisConf{
		Organization:   testGetOrganization(),
		ChannelName:    "sys-testchannel",
		ConsortiumName: "DefaultConsortium",
		Consensus: &Consensus{
			Type: "solo",
		},
		OrdererAddresses: []string{"orderer.example.com:7050"},
	}

	_, err := GetGenesisBlock(conf)
	assert.NoError(t, err, "Get genesis block error: %v", err)
}

func TestGetOrganizationJSON(t *testing.T) {
	conf := testGetOrganization()

	_, err := GetOrganizationJSON(conf)
	assert.NoError(t, err, "Get organization json error: %v", err)
}

func TestGetChannelTX(t *testing.T) {
	conf := &ChannelConf{
		Name:           "testchannel",
		ConsortiumName: "DefaultConsortium",
		Organizations: map[string]struct{}{
			"Org1": {},
			"Org2": {},
		},
	}

	_, err := GetChannelTX(conf)
	assert.NoError(t, err, "Get channel tx error: %v", err)
}

func TestGetAnchorPeerTX(t *testing.T) {
	conf := &AnchorPeerConf{
		ChannelName:      "testchannel",
		OrganizationName: "Org1",
		Peers: []*Peer{
			{
				Host: "peer0.org1.example.com",
				Port: 7051,
			},
		},
	}

	_, err := GetAnchorPeerTX(conf)
	assert.NoError(t, err, "Get anchor peer error: %v", err)
}
