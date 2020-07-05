/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package utils

import (
	"github.com/yakumioto/alkaid/internal/models"
	"github.com/yakumioto/alkaid/internal/utils/certificate"
)

func GetPkixName(org *models.Organization, commonName string) *certificate.PkixName {
	return &certificate.PkixName{
		OrgName:       org.Domain,
		CommonName:    commonName,
		Country:       org.Country,
		Province:      org.Province,
		Locality:      org.Locality,
		OrgUnit:       org.OrganizationalUnit,
		StreetAddress: org.StreetAddress,
		PostalCode:    org.PostalCode,
	}
}
