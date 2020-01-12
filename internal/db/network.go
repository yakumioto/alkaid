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
	"fmt"
	"time"

	"github.com/yakumioto/alkaid/internal/api/types"
)

type errNetwork struct {
	NetworkID string
}

type ErrNetworkExist struct {
	errNetwork
}

func (e *ErrNetworkExist) Error() string {
	return fmt.Sprintf("network already exists [network_id: %s]", e.NetworkID)
}

type ErrNetworkNotExist struct {
	errNetwork
}

func (e *ErrNetworkNotExist) Error() string {
	return fmt.Sprintf("network not exists [network_id: %s]", e.NetworkID)
}

type Network struct {
	ID                int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	NetworkID         string `xorm:"'network_id' UNIQUE INDEX NOT NULL"`
	DockerNetworkName string `xorm:"'docker_network_name' UNIQUE INDEX NOT NULL"`
	Name              string `xorm:"'name'"`
	Type              string `xorm:"'type'"`
	Description       string `xorm:"'description'"`
	CreatedAt         int64  `xorm:"'created_at'"`
	UpdatedAt         int64  `xorm:"'updated_at'"`
}

func (*Network) TableName() string {
	return "network"
}

func (n *Network) BeforeInsert() {
	n.CreatedAt = time.Now().Unix()
	n.UpdatedAt = n.CreatedAt
}

func isNetworkExist(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	return x.Get(&Network{NetworkID: id})
}

func CreateNetwork(network *Network) error {
	exist, err := isNetworkExist(network.NetworkID)
	if err != nil {
		return err
	}

	if exist {
		return &ErrNetworkExist{errNetwork{NetworkID: network.NetworkID}}
	}

	err = (*types.Network)(network).Init()
	if err != nil {
		logger.Errof("Network Init error: %v", err)
		return err
	}

	_, err = x.Insert(network)
	if err != nil {
		return err
	}

	return nil
}

func QueryNetworkByNetworkID(id string) (*types.Network, error) {
	network := &Network{NetworkID: id}

	has, err := x.Get(network)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &ErrNetworkNotExist{errNetwork{NetworkID: id}}
	}

	return (*types.Network)(network), nil
}
