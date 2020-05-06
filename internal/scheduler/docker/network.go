/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package docker

import (
	"github.com/yakumioto/alkaid/internal/api/types"
	dockervm "github.com/yakumioto/alkaid/internal/vm/docker"
)

type network struct {
	cli *dockervm.Controller
}

func newNetwork(cli *dockervm.Controller) *network {
	return &network{cli: cli}
}

func (n *network) CreateNetwork(network *types.Network) error {
	id, err := n.cli.CreateNetworkWithDockerMode(network.DockerNetworkName)
	if err != nil {
		logger.Errof("Create docker network error: %s", err)
		return err
	}

	network.DockerNetworkID = id

	logger.Debuf("Docker network id: %s", id)
	logger.Debuf("Docker network name: %s", network.DockerNetworkName)

	return nil
}

func (n *network) DeleteNetwork() error {
	return nil
}
