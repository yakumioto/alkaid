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
	ID                        string `json:"id,omitempty" gorm:"uniqueIndex"`
	Name                      string `json:"name,omitempty"`
	Domain                    string `json:"domain,omitempty" gorm:"uniqueIndex"`
	Description               string `json:"description,omitempty"`
	Type                      string `json:"type,omitempty"`
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

func (o *Organization) create() error {
	o.ResourceID = util.GenResourceID(ResourceNamespace)
	return storage.Create(o)
}
