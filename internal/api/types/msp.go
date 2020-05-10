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

type User struct {
	ID             int64  `json:"-"`
	OrganizationID string `json:"organization_id,omitempty"`

	// If the type is orderer or peer, the hostname is ({.UserID}}.{{Organization.Domain}}
	// If the type is admin or client, the name is {{.UserID}}@{{Organization.Domain}}
	UserID string `json:"user_id,omitempty" binding:"required"`
	Name   string `json:"name,omitempty" binding:"required"`

	// If type exists and is an orderer or peer, you can set SANS
	SANS []string `json:"sans,omitempty"`

	// Type is the following four orderer, peer, admin and client
	MSPType string `json:"msp_type,omitempty" binding:"required,oneof=orderer peer admin client"`

	Description     string `json:"description,omitempty"`
	NodeOUs         bool   `json:"node_o_us,omitempty"`
	PrivateKey      []byte `json:"-"`
	SignCertificate []byte `json:"sign_certificate,omitempty"`
	TLSCertificate  []byte `json:"tls_certificate,omitempty"`
	CreateAt        int64  `json:"create_at,omitempty"`
	UpdateAt        int64  `json:"update_at,omitempty"`
}

func NewUser() *User {
	return &User{}
}
