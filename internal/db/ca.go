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

type ErrCAExist struct {
	OrganizationID string
	Type           string
}

func (e *ErrCAExist) Error() string {
	return fmt.Sprintf("ca already exists [organization_id: %s, type: %s]", e.OrganizationID, e.Type)
}

type ErrCANotExist struct {
	OrganizationID string
	Type           string
}

func (e *ErrCANotExist) Error() string {
	return fmt.Sprintf("ca not exists [organization_id: %s, type: %s]", e.OrganizationID, e.Type)
}

// CA Sign CA or TLS CA
type CA struct {
	ID             int64  `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID string `xorm:"'organziation_id' UNIQUE(type~org) NOT NULL"`
	Type           string `xorm:"'type' UNIQUE(type~org) NOT NULL"`
	PrivateKey     []byte `xorm:"'private_key'"`
	Certificate    []byte `xorm:"'certificate'"`
	CreateAt       int64  `xorm:"'create_at'"`
	UpdateAt       int64  `xorm:"'update_at'"`
}

func (*CA) TableName() string {
	return "ca"
}

func (ca *CA) BeforeInsert() {
	ca.CreateAt = time.Now().Unix()
	ca.UpdateAt = ca.CreateAt
}

func isCAExist(id, typ string) (bool, error) {
	if id == "" {
		return false, nil
	}

	return x.Get(&CA{OrganizationID: id, Type: typ})
}

func CreateCA(ca *CA) error {
	exist, err := isOrganizationExist(ca.OrganizationID)
	if err != nil {
		return err
	}

	if !exist {
		return &ErrOrganizationNotExist{OrganizationID: ca.OrganizationID}
	}

	exist, err = isCAExist(ca.OrganizationID, ca.Type)
	if err != nil {
		return err
	}

	if exist {
		return &ErrCAExist{Type: ca.Type, OrganizationID: ca.OrganizationID}
	}

	_, err = x.Insert(ca)
	if err != nil {
		return err
	}

	return nil
}

func QueryCAByOrganizationIDAndType(orgID, typ string) (*types.CA, error) {
	ca := &CA{
		OrganizationID: orgID,
		Type:           typ,
	}

	has, err := x.Get(ca)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &ErrCANotExist{OrganizationID: orgID, Type: typ}
	}

	return (*types.CA)(ca), nil
}
