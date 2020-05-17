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

	"github.com/docker/docker/pkg/namesgenerator"

	"github.com/yakumioto/alkaid/internal/api/types"
)

var (
	ErrNetworkExist    = new(NetworkExistError)
	ErrNetworkNotExist = new(NetworkNotExistError)
)

type networkError struct {
	NetworkID string
}

type NetworkExistError struct {
	networkError
}

func (e *NetworkExistError) Error() string {
	return fmt.Sprintf("network already exists [network_id: %s]", e.NetworkID)
}

type NetworkNotExistError struct {
	networkError
}

func (e *NetworkNotExistError) Error() string {
	return fmt.Sprintf("network not exists [network_id: %s]", e.NetworkID)
}

type Network struct {
	ID              int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	NetworkID       string `xorm:"'network_id' UNIQUE INDEX NOT NULL"`
	Name            string `xorm:"'name'"`
	Type            string `xorm:"'type'"`
	Description     string `xorm:"'description'"`
	DockerNetworkID string `xorm:"'docker_network_id' UNIQUE INDEX NOT NULL"`
	CreatedAt       int64  `xorm:"'created_at'"`
	UpdatedAt       int64  `xorm:"'updated_at'"`
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
		return &NetworkExistError{networkError{NetworkID: network.NetworkID}}
	}

	network.DockerNetworkID = namesgenerator.GetRandomName(0)

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
		return nil, &NetworkNotExistError{networkError{NetworkID: id}}
	}

	return (*types.Network)(network), nil
}
