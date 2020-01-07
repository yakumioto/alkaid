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
	OrdererMSPType = "orderer"
	PeerMSPType    = "peer"
	AdminMSPType   = "admin"
	ClientMSPType  = "client"
)

type MSP struct {
	ID                int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	MSPID             string `xorm:"'mspid' UNIQUE(mspid~org) NOT NULL"`
	OrganizationID    string `xorm:"'organziation_id' UNIQUE(mspid~org) NOT NULL"`
	Name              string `xorm:"'name'"`
	Type              string `xorm:"'type'"` // orderer, peer, admin and client
	PrivateKey        []byte `xorm:"'private_key'" json:"-"`
	Certificate       []byte `xorm:"'certificate'"`
	SignCACertificate []byte `xorm:"'sign_ca_certificate'"`
	TLSCertificate    []byte `xorm:"'tls_certificate'"`
	CreateAt          int64  `xorm:"'create_at'"`
	UpdateAt          int64  `xorm:"'update_at'"`
}
