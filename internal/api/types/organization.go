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
	"io/ioutil"
	"os"
	"path"

	"github.com/yakumioto/alkaid/internal/config"
)

const (
	cacerts    = "cacerts"
	tlscacerts = "tlscacerts"
	admincerts = "admincerts"

	OrdererOrgType = "orderer"
	PeerOrgType    = "peer"
)

// Organization in the network
type Organization struct {
	ID             int64    `json:"-"`
	OrganizationID string   `json:"organization_id,omitempty" binding:"required"`
	Name           string   `json:"name,omitempty" binding:"required"`
	NetworkID      []string `json:"network_id,omitempty"`
	Domain         string   `json:"domain,omitempty" binding:"required,fqdn"`

	// Type value is orderer or peer
	Type string `json:"type,omitempty" binding:"required,oneof=orderer peer"`

	Description string `json:"description,omitempty"`

	// The following fields are the fields that generate the certificate
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	Locality           string `json:"locality,omitempty"`
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
	StreetAddress      string `json:"street_address,omitempty"`
	PostalCode         string `json:"postal_code,omitempty"`
	MSPDir             string `json:"msp_dir,omitempty"`
	CreateAt           int64  `json:"create_at,omitempty"`
	UpdateAt           int64  `json:"update_at,omitempty"`
}

// NewOrganization Default parameter
func NewOrganization() *Organization {
	return &Organization{
		Country:    "China",
		Province:   "Beijing",
		Locality:   "Beijing",
		PostalCode: "100000",
	}
}

// CreateMSPDir The configtxgen tool needs to be used
func (o *Organization) CreateMSPDir(signca, tlsca *CA) error {
	dirs := make(map[string]string)

	rootDir := path.Join(path.Clean(config.FileSystemPath), o.OrganizationID, "msp")
	o.MSPDir = rootDir

	dirs[admincerts] = path.Join(rootDir, admincerts)
	dirs[cacerts] = path.Join(rootDir, cacerts)
	dirs[tlscacerts] = path.Join(rootDir, tlscacerts)

	for base, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		switch base {
		case cacerts:
			if err := ioutil.WriteFile(path.Join(dir, "cert.pem"), signca.Certificate, 0644); err != nil {
				return err
			}
		case tlscacerts:
			if err := ioutil.WriteFile(path.Join(dir, "cert.pem"), tlsca.Certificate, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
