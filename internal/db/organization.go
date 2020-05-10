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

type ErrOrganizationExist struct {
	OrganizationID string
}

func (e *ErrOrganizationExist) Error() string {
	return fmt.Sprintf("organization already exists [organization_id: %s]", e.OrganizationID)
}

type ErrOrganizationNotExist struct {
	OrganizationID string
}

func (e *ErrOrganizationNotExist) Error() string {
	return fmt.Sprintf("organization not exists [organization_id: %s]", e.OrganizationID)
}

// Organization in the network
type Organization struct {
	ID                 int64    `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID     string   `xorm:"'organziation_id' UNIQUE NOT NULL"`
	Name               string   `xorm:"'name'"`
	NetworkID          []string `xorm:"'network_id'"`
	Domain             string   `xorm:"'domain' UNIQUE NOT NULL"`
	Description        string   `xorm:"'description'"`
	Type               string   `xorm:"'type'"` // orderer or peer
	Country            string   `xorm:"'country'"`
	Province           string   `xorm:"'province'"`
	Locality           string   `xorm:"'locality'"`
	OrganizationalUnit string   `xorm:"'organizational_unit'"`
	StreetAddress      string   `xorm:"'street_address'"`
	PostalCode         string   `xorm:"'postal_code'"`
	SignCAPrivateKey   []byte   `xorm:"'sign_ca_private_key'"`
	TLSCAPrivateKey    []byte   `xorm:"'tlsca_private_key'"`
	SignCACertificate  []byte   `xorm:"'sign_ca_certificate'"`
	TLSCACertificate   []byte   `xorm:"'tlsca_certificate'"`
	CreateAt           int64    `xorm:"'create_at'"`
	UpdateAt           int64    `xorm:"'update_at'"`
}

func (org *Organization) BeforeInsert() {
	org.CreateAt = time.Now().Unix()
	org.UpdateAt = org.CreateAt
}

func isOrganizationExist(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	return x.Get(&Organization{OrganizationID: id})
}

func CreateOrganization(org *Organization) error {
	exist, err := isOrganizationExist(org.OrganizationID)
	if err != nil {
		return err
	}

	if exist {
		return &ErrOrganizationExist{OrganizationID: org.OrganizationID}
	}

	_, err = x.Insert(org)
	if err != nil {
		return err
	}

	return nil
}

func QueryOrganizationByOrgID(orgID string) (*types.Organization, error) {
	org := &Organization{
		OrganizationID: orgID,
	}

	has, err := x.Get(org)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &ErrOrganizationNotExist{OrganizationID: orgID}
	}

	return (*types.Organization)(org), nil
}
