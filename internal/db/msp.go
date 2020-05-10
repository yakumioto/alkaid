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

type ErrUser struct {
	OrganizationID string
	UserID         string
}

type ErrUserExist struct {
	ErrUser
}

func (e *ErrUserExist) Error() string {
	return fmt.Sprintf("user already exists [organization_id: %s, user_id: %s]", e.OrganizationID, e.UserID)
}

type ErrUserNotExist struct {
	ErrUser
}

func (e *ErrUserNotExist) Error() string {
	return fmt.Sprintf("user not exists [organization_id: %s, user_id: %s]", e.OrganizationID, e.UserID)
}

type User struct {
	ID              int64    `xorm:"'id' PRIMARY KEY AUTOINCREMENT NOT NULL"`
	OrganizationID  string   `xorm:"'organziation_id' UNIQUE(userid~orgid) NOT NULL"`
	UserID          string   `xorm:"'user_id' UNIQUE(userid~orgid) NOT NULL"`
	Name            string   `xorm:"'name'"`
	SANS            []string `xorm:"'sans'"`
	MSPType         string   `xorm:"'msp_type'"`
	Description     string   `xorm:"'description'"`
	NodeOUs         bool     `xorm:"'node_ous'"`
	PrivateKey      []byte   `xorm:"'private_key'"`
	SignCertificate []byte   `xorm:"'sign_ca_certificate'"`
	TLSCertificate  []byte   `xorm:"'tls_certificate'"`
	CreateAt        int64    `xorm:"'create_at'"`
	UpdateAt        int64    `xorm:"'update_at'"`
}

func (*User) TableName() string {
	return "user"
}

func (m *User) BeforeInsert() {
	m.CreateAt = time.Now().Unix()
	m.UpdateAt = m.CreateAt
}
func isMSPExist(id, orgID string) (bool, error) {
	if id == "" || orgID == "" {
		return false, nil
	}

	return x.Get(&User{UserID: id, OrganizationID: orgID})
}

func CreateMSP(msp *User) error {
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
		return &ErrUserExist{ErrUser{OrganizationID: msp.OrganizationID, UserID: msp.UserID}}
	}

	_, err = x.Insert(msp)
	if err != nil {
		return err
	}

	return nil
}

func QueryMSPByOrganizationIDAndUserID(orgID, userID string) (*types.User, error) {
	user := &User{
		OrganizationID: orgID,
		UserID:         userID,
	}

	has, err := x.Get(user)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, &ErrUserNotExist{ErrUser{
			OrganizationID: orgID,
			UserID:         userID,
		}}
	}

	return (*types.User)(user), nil
}
