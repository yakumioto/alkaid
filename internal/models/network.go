/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package models

type Network struct {
	ID                     string   `json:"id,omitempty" xorm:"'id' index notnull"`
	NetworkID              string   `json:"docker_network_id,omitempty" xorm:"'docker_network_id'"`
	Name                   string   `json:"name,omitempty" xorm:"'name'"`
	Type                   string   `json:"type,omitempty" xorm:"'type'"`
	ConsensusType          string   `json:"consensus_type,omitempty" xorm:"'consensus_type'"`
	OrdererManagement      bool     `json:"orderer_management,omitempty" xorm:"'orderer_management'"`
	OrdererOrganizationIDs []string `json:"orderer_organization_ids,omitempty" xorm:"'orderer_organization_ids'"`
	PeerOrganizationIDs    []string `json:"peer_organization_ids,omitempty" xorm:"'peer_organization_ids'"`

	Description string `json:"description,omitempty" xorm:"'description'"`
	CreatedAt   int64  `json:"created_at,omitempty" xorm:"'created_at' created"`
	UpdatedAt   int64  `json:"updated_at,omitempty" xorm:"'updated_at' updated"`
}

func NewNetwork() *Network {
	return &Network{
		OrdererManagement: true,
	}
}

func (n *Network) TableName() string {
	return "networks"
}

func (n *Network) Create() error {
	return db.Create(n)
}

func (n *Network) Update(id string) error {
	return db.Update(&Network{ID: id}, n)
}

func (n *Network) Get(id string) error {
	n.ID = id
	return db.Get("", n)
}

func (n *Network) Delete(id string) error {
	return db.Delete(&Network{ID: id})
}

type Networks []*Network

func (o *Networks) Query(query ...interface{}) error {
	return db.Query(o, query...)
}
