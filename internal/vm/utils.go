/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package vm

import (
	"errors"

	"github.com/yakumioto/alkaid/internal/api/types"
	"github.com/yakumioto/alkaid/internal/vm/docker"
	"github.com/yakumioto/alkaid/third_party/github.com/moby/moby/pkg/namesgenerator"
)

func NetworkInitialize(network *types.Network) error {
	switch machine := vm.(type) {
	case *docker.Controller:
		if network.Type == types.DockerNetworkType {
			name := namesgenerator.GetRandomName(0)
			id, err := machine.CreateNetworkWithDockerMode(name)
			if err != nil {
				logger.Errof("Create docker network error: %s", err)
				return err
			}
			network.DockerNetworkID = id
			network.DockerNetworkName = name

			logger.Debuf("docker network id: %s", id)
			logger.Debuf("docker network name: %s", name)
		}

	default:
		return errors.New("no virtual machine interface implementation")
	}

	return nil
}
