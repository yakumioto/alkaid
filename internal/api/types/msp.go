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
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
)

const (
	OrdererMSPType = "orderer"
	PeerMSPType    = "peer"
	AdminMSPType   = "admin"
	ClientMSPType  = "client"
)

type MSP struct {
	ID             int64  `json:"-"`
	OrganizationID string `json:"organization_id,omitempty"`
	UserID         string `json:"user_id,omitempty" binding:"required"`
	// If the type is orderer or peer, the default Hostname is ({.Name}}.{{Organization.Domain}}
	// If the type is admin or client, the default name is {{.Name}}@{{Organization.Domain}}
	Name string `json:"name,omitempty" binding:"required"`

	// If type exists and is an orderer or peer, you can set SANS
	SANS []string `json:"sans,omitempty"`

	// Type is the following four orderer, peer, admin and client
	Type string `json:"type,omitempty" binding:"required,oneof=orderer peer admin client"`

	Description     string `json:"description,omitempty"`
	NodeOUs         bool   `json:"node_o_us,omitempty"`
	PrivateKey      []byte `json:"-"`
	SignCertificate []byte `json:"sign_certificate,omitempty"`
	TLSCertificate  []byte `json:"tls_certificate,omitempty"`
	CreateAt        int64  `json:"create_at,omitempty"`
	UpdateAt        int64  `json:"update_at,omitempty"`
}

func NewMSP() *MSP {
	return &MSP{}
}

func (m *MSP) Initialize(org *Organization, signCA, tlsCA *CA) error {
	priv, err := crypto.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("generate private key error: %v", err)
	}

	signCert, err := signCA.SignCertificate(org, []string{m.Type}, nil, &priv.PublicKey)
	if err != nil {
		return fmt.Errorf("sign certificate error: %v", err)
	}

	tlsCert, err := tlsCA.SignCertificate(org, []string{m.Type}, m.SANS, &priv.PublicKey)
	if err != nil {
		return fmt.Errorf("sign certificate error: %v", err)
	}

	privBytes, err := crypto.PrivateKeyExport(priv)
	if err != nil {
		return fmt.Errorf("private key export error: %v", err)
	}

	m.PrivateKey = privBytes
	m.SignCertificate = crypto.X509Export(signCert)
	m.TLSCertificate = crypto.X509Export(tlsCert)

	return nil
}
