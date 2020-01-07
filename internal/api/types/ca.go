/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package types

const (
	SignCAType = "sign_ca"
	TLSCAType  = "tls_ca"
)

// CA Sign CA or TLS CA
type CA struct {
	ID             int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID int64  `xorm:"'organziation_id' UNIQUE(type~org) NOT NULL"`
	Type           string `xorm:"'type' UNIQUE(type~org) NOT NULL"`
	PrivateKey     []byte `xorm:"'private_key'" json:"-"`
	Certificate    []byte `xorm:"'certificate'"`
	CreateAt       int64  `xorm:"'create_at'"`
	UpdateAt       int64  `xorm:"'update_at'"`
}
