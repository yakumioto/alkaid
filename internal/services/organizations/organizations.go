/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package organizations

import (
	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/util"
)

const ResourceNamespace = "Organization"

// Organization 组织，组织中包含了加密后的 Sign CA，TLS CA 密钥。
// 所以在创建组织时需要填入一个交易密码，此密码用来加解密上述的两个 CA 密钥。
type Organization struct {
	ResourceID                string `json:"resourceID,omitempty" gorm:"primaryKey"`
	OrganizationID            string `json:"organizationId,omitempty" gorm:"uniqueIndex"`
	Name                      string `json:"name,omitempty"`
	Domain                    string `json:"domain,omitempty" gorm:"uniqueIndex"`
	Description               string `json:"description,omitempty"`
	Country                   string `json:"country,omitempty"`
	Province                  string `json:"province,omitempty"`
	Locality                  string `json:"locality,omitempty"`
	OrganizationalUnit        string `json:"organizationalUnit,omitempty"`
	StreetAddress             string `json:"streetAddress,omitempty"`
	PostalCode                string `json:"postalCode,omitempty"`
	ProtectedSignCAPrivateKey string `json:"protectedSignCAPrivateKey,omitempty"`
	ProtectedTLSCAPrivateKey  string `json:"protectedTlsCAPrivateKey,omitempty"`
	SignCACertificate         string `json:"signCACertificate,omitempty"`
	TlsCACertificate          string `json:"tlsCACertificate,omitempty"`
	CreatedAt                 int64  `json:"createdAt,omitempty"`
	UpdatedAt                 int64  `json:"updatedAt,omitempty"`
}

func newOrganizationByCreateRequest(req *CreateRequest) *Organization {
	org := &Organization{
		OrganizationID:     req.OrganizationID,
		Name:               req.Name,
		Domain:             req.Domain,
		Description:        req.Description,
		Country:            req.Country,
		Province:           req.Province,
		Locality:           req.Locality,
		OrganizationalUnit: req.OrganizationalUnit,
		StreetAddress:      req.StreetAddress,
		PostalCode:         req.PostalCode,
	}

	org.SetCountry(org.Country)
	org.SetProvince(org.Province)
	org.SetLocality(org.Locality)
	org.SetOrganizationalUnit(org.OrganizationalUnit)

	return org
}

func (o *Organization) Create() error {
	o.ResourceID = util.GenResourceID(ResourceNamespace)
	o.SetCountry(o.Country)
	o.SetProvince(o.Province)
	o.SetLocality(o.Locality)
	o.SetOrganizationalUnit(o.OrganizationalUnit)
	return storage.Create(o)
}

func (o *Organization) SetCountry(country string) {
	if country != "" {
		o.Country = country
	}

	o.Country = "China"
}

func (o *Organization) SetProvince(province string) {
	if province != "" {
		o.Province = province
	}

	o.Province = "Beijing"
}

func (o *Organization) SetLocality(locality string) {
	if locality != "" {
		o.Locality = locality
	}

	o.Locality = "Beijing"
}

func (o *Organization) SetOrganizationalUnit(organizationalUnit string) {
	if organizationalUnit != "" {
		o.OrganizationalUnit = organizationalUnit
	}

	o.OrganizationalUnit = "Alkaid"
}

func FindOrganizationByID(id string) (*Organization, error) {
	return nil, nil
}
