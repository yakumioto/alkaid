/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package db

type Container struct {
	ID             int64
	OrganizationID string
	ContainerID    string
	ContainerName  string
	NetworkID      string
	Type           string
	ImageName      string
	ImageTag       string
	Volumes        []string
	Status         bool
	CreateAt       int64
	UpdateAt       int64
}
