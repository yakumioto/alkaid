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

type ErrMSP struct {
	OrganizationID string
	UserID         string
}

type ErrMSPExist struct {
	ErrMSP
}

func (e *ErrMSPExist) Error() string {
	return fmt.Sprintf("msp already exists [organization_id: %s, user_id: %s]", e.OrganizationID, e.UserID)
}

type ErrMSPNotExist struct {
	ErrMSP
}

func (e *ErrMSPNotExist) Error() string {
	return fmt.Sprintf("msp not exists [organization_id: %s, user_id: %s]", e.OrganizationID, e.UserID)
}

type MSP struct {
	ID              int64    `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID  string   `xorm:"'organziation_id' UNIQUE(userid~orgid) NOT NULL"`
	UserID          string   `xorm:"'user_id' UNIQUE(userid~orgid) NOT NULL"`
	Name            string   `xorm:"'name'"`
	SANS            []string `xorm:"'sans'"`
	Type            string   `xorm:"'type'"`
	Description     string   `xorm:"'description'"`
	NodeOUs         bool     `xorm:"'node_ous'"`
	PrivateKey      []byte   `xorm:"'private_key'"`
	SignCertificate []byte   `xorm:"'sign_ca_certificate'"`
	TLSCertificate  []byte   `xorm:"'tls_certificate'"`
	CreateAt        int64    `xorm:"'create_at'"`
	UpdateAt        int64    `xorm:"'update_at'"`
}

func (*MSP) TableName() string {
	return "msp"
}

func (m *MSP) BeforeInsert() {
	m.CreateAt = time.Now().Unix()
	m.UpdateAt = m.CreateAt
}
func isMSPExist(id, orgID string) (bool, error) {
	if id == "" || orgID == "" {
		return false, nil
	}

	return x.Get(&MSP{UserID: id, OrganizationID: orgID})
}

func CreateMSP(msp *MSP) error {
	exist, err := isOrganizationExist(msp.OrganizationID)
	if err != nil {
		return err
	}

	if !exist {
		return &ErrOrganizationNotExist{OrganizationID: msp.OrganizationID}
	}

	exist, err = isMSPExist(msp.UserID, msp.OrganizationID)
	if err != nil {
		return err
	}

	if exist {
		return &ErrMSPExist{ErrMSP{OrganizationID: msp.OrganizationID, UserID: msp.UserID}}
	}

	_, err = x.Insert(msp)
	if err != nil {
		return err
	}

	return nil
}

func QueryMSPByOrganizationIDAndUserID(orgID, userID string) (*types.MSP, error) {
	msp := &MSP{
		OrganizationID: orgID,
		UserID:         userID,
	}

	has, err := x.Get(msp)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &ErrMSPNotExist{ErrMSP{
			OrganizationID: orgID,
			UserID:         userID,
		}}
	}

	return (*types.MSP)(msp), nil
}
