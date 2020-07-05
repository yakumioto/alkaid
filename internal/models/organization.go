/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package models

type Organization struct {
	ID     string `json:"id,omitempty" binding:"required,alphanum" xorm:"'id' index notnull"`
	Name   string `json:"name,omitempty" binding:"required" xorm:"'name'"`
	Domain string `json:"domain,omitempty" binding:"required,fqdn" xorm:"'domain' unique notnull"`
	Type   string `json:"type,omitempty" binding:"required,oneof=Orderer Peer" xorm:"'type'"`

	// The following fields are the fields that generate the certificate
	Country            string `json:"country,omitempty" xorm:"'country'"`
	Province           string `json:"province,omitempty" xorm:"'province'"`
	Locality           string `json:"locality,omitempty" xorm:"'locality'"`
	OrganizationalUnit string `json:"organizational_unit,omitempty" xorm:"'organizational_unit'"`
	StreetAddress      string `json:"street_address,omitempty" xorm:"'street_address'"`
	PostalCode         string `json:"postal_code,omitempty" xorm:"'postal_code'"`

	CAPrivateKey     []byte `json:"ca_private_key,omitempty" xorm:"'ca_private_key'"`
	TLSCAPrivateKey  []byte `json:"tlsca_private_key,omitempty" xorm:"'tlsca_private_key'"`
	CACertificate    []byte `json:"ca_certificate,omitempty" xorm:"'ca_certificate'"`
	TLSCACertificate []byte `json:"tlsca_certificate,omitempty" xorm:"'tlsca_certificate'"`

	Description string `json:"description,omitempty" xorm:"'description'"`
	CreatedAt   int64  `json:"created_at,omitempty" xorm:"'created_at' created"`
	UpdatedAt   int64  `json:"updated_at,omitempty" xorm:"'updated_at' updated"`
}

func NewOrganization() *Organization {
	return &Organization{
		Country:    "China",
		Province:   "Beijing",
		Locality:   "Beijing",
		PostalCode: "100000",
	}
}

func (o *Organization) TableName() string {
	return "organizations"
}

func (o *Organization) Create() error {
	return db.Create(o)
}

func (o *Organization) Update(id string) error {
	return db.Update(&Organization{ID: id}, o)
}

func (o *Organization) Get(id string) error {
	o.ID = id
	return db.Get("", o)
}

func (o *Organization) Delete(id string) error {
	return db.Delete(&Organization{ID: id})
}

type Organizations []*Organization

func (o *Organizations) Query(query ...interface{}) error {
	return db.Query(o, query...)
}
