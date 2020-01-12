/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package types

import (
	"github.com/yakumioto/alkaid/internal/vm"
	"github.com/yakumioto/alkaid/internal/vm/docker"
	"github.com/yakumioto/alkaid/third_party/github.com/moby/moby/pkg/namesgenerator"
)

const (
	DockerNetworkType      = "docker"
	DockerSwarmNetworkType = "docker_swarm"
	KubernetesNetworkType  = "kubernetes"
)

type Network struct {
	ID                int64  `json:"id,omitempty"`
	NetworkID         string `json:"network_id,omitempty" binding:"required"`
	DockerNetworkName string `json:"-"`
	Name              string `json:"name,omitempty"`
	Type              string `json:"type,omitempty" binding:"required"`
	Description       string `json:"description,omitempty"`
	CreatedAt         int64  `json:"created_at,omitempty"`
	UpdatedAt         int64  `json:"updated_at,omitempty"`
}

func NewNetwork() *Network {
	return &Network{}
}

func (n *Network) Init() error {
	switch machine := vm.VM.(type) {
	case *docker.Controller:
		if n.Type == DockerNetworkType {
			name := namesgenerator.GetRandomName(0)
			if err := machine.CreateNetworkWithDockerMode(name); err != nil {
				return err
			}
			n.DockerNetworkName = name
		}
	}

	return nil
}
