/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package models

type User struct {
	ID             string `json:"id,omitempty" binding:"required,alphanum" xorm:"'id' index notnull"`
	OrganizationID string `json:"organization_id,omitempty" binding:"required" xorm:"'organization_id'"`
	Name           string `json:"name,omitempty" binding:"required" xorm:"'name'"`
	Password       string `json:"password,omitempty" binding:"required" xorm:"'password'"`
	MSPType        string `json:"msp_type,omitempty" binding:"required,oneof=Admin Client Orderer Peer" xorm:"'msp_type'"`
	NodeOUs        string `json:"node_o_us,omitempty" binding:"required" xorm:"'node_o_us'"`

	SANs            []string `json:"sans,omitempty" xorm:"'sans'"`
	SignPrivateKey  []byte   `json:"sign_private_key,omitempty"`
	TLSPrivateKey   []byte   `json:"tls_private_key,omitempty" xorm:"'tls_private_key'"`
	SignCertificate []byte   `json:"sign_certificate,omitempty" xorm:"'sign_certificate'"`
	TLSCertificate  []byte   `json:"tls_certificate,omitempty" xorm:"'tls_certificate'"`

	Description string `json:"description,omitempty" xorm:"'description'"`
	CreatedAt   int64  `json:"created_at,omitempty" xorm:"'created_at' created"`
	UpdatedAt   int64  `json:"updated_at,omitempty" xorm:"'updated_at' updated"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) Create() error {
	return db.Create(u)
}

func (u *User) Update(id string) error {
	return db.Update(&User{ID: id}, u)
}

func (u *User) Get(id string) error {
	u.ID = id
	return db.Get("", u)
}

func (u *User) Delete(id string) error {
	return db.Delete(&User{ID: id})
}

type Users []*User

func (u *Users) Query(query ...interface{}) error {
	return db.Query(u, query...)
}
