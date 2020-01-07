/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package types

const (
	OrdererOrgType = "orderer"
	PeerOrgType    = "peer"
)

// Organization in the network
type Organization struct {
	ID                 int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID     string `xorm:"'organziation_id' UNIQUE((org~networkid)) NOT NULL"`
	NetworkID          int64  `xorm:"UNIQUE(org~networkid) NOT NULL"`
	Domain             string `xorm:"UNIQUE NOT NULL"`
	Type               string `xorm:"'type'"` // orderer or peer
	Name               string `xorm:"'name'"`
	Description        string `xorm:"'description'"`
	Country            string `xorm:"'country'"`
	Province           string `xorm:"'province'"`
	Locality           string `xorm:"'locality'"`
	OrganizationalUnit string `xorm:"'organizational_unit'"`
	StreetAddress      string `xorm:"'street_address'"`
	PostalCode         string `xorm:"'postal_code'"`
	CreateAt           int64  `xorm:"'create_at'"`
	UpdateAt           int64  `xorm:"'update_at'"`
}
